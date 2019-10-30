package service

import (
	"context"
	"net"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	protoempty "github.com/gogo/protobuf/types"
	"github.com/jinzhu/copier"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/auth"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-users/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type RpcServerOptions struct {
	Addr               string
	AuthTokenSecret    string
	AuthRecoverySecret string

	Logger   *logrus.Entry
	DS       *datastore.Datastore
	Accounts accountsv1.AccountServiceClient
	EB       *EventBus
}

type RpcServer struct {
	addr               string
	authTokenSecret    string
	authRecoverySecret string

	grpc          *grpc.Server
	listen        net.Listener
	logger        *logrus.Entry
	ds            *datastore.Datastore
	eb            *EventBus
	accounts      accountsv1.AccountServiceClient
	notifications *NotificationClient
	validator     *requestValidator
}

func NewRpcServer(opts *RpcServerOptions) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	nblogger := opts.Logger.WithField("system", "notification")
	nc, err := NewNotificationClient(opts.EB, nblogger)
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		addr:               opts.Addr,
		authTokenSecret:    opts.AuthTokenSecret,
		authRecoverySecret: opts.AuthRecoverySecret,
		grpc:               grpcServer,
		listen:             listen,
		logger:             opts.Logger,
		ds:                 opts.DS,
		eb:                 opts.EB,
		accounts:           opts.Accounts,
		notifications:      nc,
		validator:          newRequestValidator(),
	}

	v1.RegisterUserServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.TokenResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)
	span.SetTag("name", req.Name)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.Register(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		s.logger.Error(err)
		if err == datastore.ErrUserAlreadyExists {
			respErr := &rpc.MultiValidationError{
				Errors: []*rpc.ValidationError{
					&rpc.ValidationError{
						Field:   "email",
						Message: "Email is already registered",
					},
				},
			}
			return nil, rpc.NewRpcValidationError(respErr)
		}

		return nil, rpc.ErrRpcInternal
	}

	token, err := s.createToken(ctx, user, v1.TokenTypeRegular)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.UpdateAuthToken(ctx, user, token); err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	accReq := &accountsv1.AccountRequest{OwnerId: user.Id}
	if err = s.eb.CreateUserAccount(span, accReq); err != nil {
		s.logger.Errorf("failed to create account via eventbus: %s", err)
	}

	if err = s.notifications.SendEmailWaitlisted(ctx, user); err != nil {
		s.logger.WithField("failed to send whitelisted email to user id", user.Id).Error(err)
	}

	return &v1.TokenResponse{
		Token: token,
	}, nil
}

func (s *RpcServer) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.TokenResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == datastore.ErrUserNotFound {
			return nil, rpc.ErrRpcUnauthenticated
		}

		return nil, rpc.ErrRpcInternal
	}

	if !checkPasswordHash(ctx, req.Password, user.Password) {
		return nil, rpc.ErrRpcUnauthenticated
	}

	token, err := s.createToken(ctx, user, v1.TokenTypeRegular)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.UpdateAuthToken(ctx, user, token); err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return &v1.TokenResponse{
		Token: user.Token,
	}, nil
}

func (s *RpcServer) Logout(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if err = s.ds.User.ResetAuthToken(user); err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) StartRecovery(ctx context.Context, req *v1.StartRecoveryUserRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == datastore.ErrUserNotFound {
			return nil, rpc.ErrRpcBadRequest
		}

		s.logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	token := newRecoveryToken(req.Email, 12*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))

	if err = s.notifications.SendEmailRecovery(ctx, user, token); err != nil {
		s.logger.WithField("failed to send recovery email to user id", user.Id).Error(err)
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Recover(ctx context.Context, req *v1.RecoverUserRequest) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := verifyRecoveryToken(ctx, req.Token, s.ds.User.GetByEmail, []byte(s.authRecoverySecret))
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == datastore.ErrUserNotFound {
			return nil, rpc.ErrRpcBadRequest
		}

		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.ResetPassword(ctx, user, req.Password); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) ResetPassword(ctx context.Context, req *v1.ResetPasswordUserRequest) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.ds.User.ResetPassword(ctx, user, req.Password); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Get(ctx context.Context, req *protoempty.Empty) (*v1.UserProfile, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	accountProfile, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: user.Id})
	if err != nil {
		s.logger.Errorf("failed to get account profile: %s", err)
	} else {
		userProfile.Account = accountProfile
	}

	return userProfile, nil
}

func (s *RpcServer) Whitelist(ctx context.Context, req *protoempty.Empty) (*v1.WhitelistResponse, error) {
	_ = opentracing.SpanFromContext(ctx)

	accounts, err := s.accounts.List(ctx, new(protoempty.Empty))
	if err != nil {
		s.logger.Errorf("failed to get whitelist: %s", err)
		return nil, err
	}

	items := make([]string, 0)

	for _, a := range accounts.Items {
		items = append(items, a.Address)
	}

	return &v1.WhitelistResponse{
		Items: items,
	}, nil
}

func (s *RpcServer) LookupByAddress(ctx context.Context, req *v1.LookupByAddressRequest) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	aReq := &accountsv1.Address{Address: req.Address}
	_, err := s.accounts.GetByAddress(ctx, aReq)
	if err != nil {
		s.logger.Errorf("failed to look up address: %s", err)
		return nil, rpc.ErrRpcNotFound

	}

	return new(protoempty.Empty), nil
}

func (s *RpcServer) Activate(ctx context.Context, req *v1.UserRequest) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	requester, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if requester.Role < v1.UserRoleManager {
		return nil, rpc.ErrRpcPermissionDenied
	}

	user, err := s.ds.User.Get(req.Id)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == datastore.ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.ErrRpcInternal
	}

	if err := s.ds.User.Activate(req.Id); err != nil {
		s.logger.Errorf("failed to activate user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.notifications.SendEmailWelcome(ctx, user); err != nil {
		s.logger.WithField("failed to send welcome email to user id", req.Id).Error(err)
	}

	return new(protoempty.Empty), nil
}

func (s *RpcServer) Key(ctx context.Context, req *v1.UserRequest) (*v1.KeyResponse, error) {
	_ = opentracing.SpanFromContext(ctx)

	requester, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if requester.Role < v1.UserRoleManager {
		return nil, rpc.ErrRpcPermissionDenied
	}

	aReq := &accountsv1.AccountRequest{OwnerId: req.Id}
	key, err := s.accounts.Key(ctx, aReq)
	if err != nil {
		s.logger.Errorf("failed to get key: %s", err)
		return nil, rpc.ErrRpcNotFound

	}

	keyResp := &v1.KeyResponse{
		Key: key.Key,
	}

	return keyResp, nil
}

func (s *RpcServer) ListApiTokens(ctx context.Context, req *protoempty.Empty) (*v1.UserApiListResponse, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	tokens, err := s.ds.Token.ListByUser(ctx, user.Id)
	if err != nil {
		return nil, rpc.ErrRpcInternal
	}

	tokensResponse := []*v1.UserApiTokenResponse{}
	if err = copier.Copy(&tokensResponse, tokens); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &v1.UserApiListResponse{
		Items: tokensResponse,
	}, nil
}

func (s *RpcServer) CreateApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*v1.CreateUserApiTokenResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("name", req.Name)

	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	token, err := s.createToken(ctx, user, v1.TokenTypeAPI)
	if err != nil {
		s.logger.Errorf("failed to create api token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	apiToken, err := s.ds.Token.Create(ctx, user.Id, req.Name, token)
	if err != nil {
		s.logger.Errorf("failed to create api token record: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &v1.CreateUserApiTokenResponse{
		Id:    apiToken.Id,
		Name:  apiToken.Name,
		Token: apiToken.Token,
	}, nil
}

func (s *RpcServer) DeleteApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)

	_, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.ds.Token.Delete(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("failed to delete api token record: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return new(protoempty.Empty), nil
}

func (s *RpcServer) StartWithdraw(ctx context.Context, req *v1.StartWithdrawRequest) (*v1.WithdrawResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("address", req.Address)
	span.SetTag("amount", req.Amount)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	transfer, err := s.ds.Transfer.Create(ctx, user.Id, req.Address, req.Amount)
	if err != nil {
		s.logger.WithError(err).Error("failed to create transfer")
		return nil, rpc.ErrRpcInternal
	}

	if err = s.notifications.SendWithdrawTransfer(ctx, user, transfer); err != nil {
		s.logger.WithError(err).Error("failed to send withdraw transfer email")
	}

	return &v1.WithdrawResponse{
		TransferId: transfer.Id,
	}, nil
}

func (s *RpcServer) Withdraw(ctx context.Context, req *v1.WithdrawRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.TransferId)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	transfer, err := s.ds.Transfer.Get(ctx, req.TransferId)
	if err != nil {
		s.logger.WithError(err).Error("failed to get transfer")
		return nil, rpc.ErrRpcInternal
	}

	if transfer.Pin != req.Pin {
		s.logger.Error("failed with incorrect pin")
		return nil, rpc.ErrRpcBadRequest
	}

	if transfer.ExpiresAt.Before(time.Now()) {
		s.logger.Error("failed with expired transfer")
		return nil, rpc.ErrRpcBadRequest
	}

	_, err = s.accounts.Withdraw(ctx,
		&accountsv1.WithdrawRequest{
			OwnerId:    user.Id,
			TransferId: transfer.Id,
		},
	)
	if err != nil {
		s.logger.WithError(err).Error("failed to withdraw")
		return nil, rpc.ErrRpcInternal
	}

	return new(protoempty.Empty), nil
}

func (s *RpcServer) createToken(ctx context.Context, user *v1.User, tokenType v1.TokenType) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "createToken")
	defer span.Finish()

	span.SetTag("id", user.Id)
	span.SetTag("email", user.Email)
	span.SetTag("token_type", tokenType)

	claims := auth.ExtendedClaims{
		Type: auth.TokenType(tokenType),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Id,
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.authTokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *RpcServer) authenticate(ctx context.Context) (*v1.User, context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, err := auth.AuthFromContext(ctx)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	if s.getTokenType(ctx) == auth.TokenType(v1.TokenTypeAPI) {
		return nil, nil, rpc.ErrRpcPermissionDenied
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	user, err := s.ds.User.Get(userID)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	if user.Token == "" {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	return user, ctx, nil
}

func (s *RpcServer) getTokenType(ctx context.Context) auth.TokenType {
	tokenType, ok := auth.TypeFromContext(ctx)
	if !ok {
		return auth.TokenType(v1.TokenTypeRegular)
	}

	return tokenType
}
