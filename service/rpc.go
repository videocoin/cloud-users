package service

import (
	"context"
	"math/big"
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
	ds "github.com/videocoin/cloud-users/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RPCServerOptions struct {
	Addr               string
	AuthTokenSecret    string
	AuthRecoverySecret string

	Logger   *logrus.Entry
	DS       *datastore.Datastore
	Accounts accountsv1.AccountServiceClient
	EB       *EventBus
}

type RPCServer struct {
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

func NewRPCServer(opts *RPCServerOptions) (*RPCServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	gRPCServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	nblogger := opts.Logger.WithField("system", "notification")
	nc, err := NewNotificationClient(opts.EB, nblogger)
	if err != nil {
		return nil, err
	}
	validator, err := newRequestValidator()
	if err != nil {
		return nil, err
	}
	RPCServer := &RPCServer{
		addr:               opts.Addr,
		authTokenSecret:    opts.AuthTokenSecret,
		authRecoverySecret: opts.AuthRecoverySecret,
		grpc:               gRPCServer,
		listen:             listen,
		logger:             opts.Logger,
		ds:                 opts.DS,
		eb:                 opts.EB,
		accounts:           opts.Accounts,
		notifications:      nc,
		validator:          validator,
	}

	v1.RegisterUserServiceServer(gRPCServer, RPCServer)
	reflection.Register(gRPCServer)

	return RPCServer, nil
}

func (s *RPCServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RPCServer) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.TokenResponse, error) {
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
					{
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

	accReq := &accountsv1.AccountRequest{OwnerId: user.ID}
	if err = s.eb.CreateUserAccount(span, accReq); err != nil {
		s.logger.Errorf("failed to create account via eventbus: %s", err)
	}

	confirmToken := newRecoveryToken(req.Email, 1*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))

	if err = s.notifications.SendEmailConfirmation(ctx, user, confirmToken); err != nil {
		s.logger.WithField("failed to send confirm email to user id", user.ID).Error(err)
	}

	return &v1.TokenResponse{
		Token: token,
	}, nil
}

func (s *RPCServer) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.TokenResponse, error) {
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

func (s *RPCServer) Logout(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
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

func (s *RPCServer) StartRecovery(ctx context.Context, req *v1.StartRecoveryUserRequest) (*protoempty.Empty, error) {
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
		s.logger.WithField("failed to send recovery email to user id", user.ID).Error(err)
	}

	return &protoempty.Empty{}, nil
}

func (s *RPCServer) Recover(ctx context.Context, req *v1.RecoverUserRequest) (*protoempty.Empty, error) {
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

func (s *RPCServer) StartConfirmation(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	confirmToken := newRecoveryToken(user.Email, 1*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))

	if err = s.notifications.SendEmailConfirmation(ctx, user, confirmToken); err != nil {
		s.logger.WithField("failed to send confirmation email to user id", user.ID).Error(err)
	}

	return &protoempty.Empty{}, nil
}

func (s *RPCServer) Confirm(ctx context.Context, req *v1.ConfirmUserRequest) (*protoempty.Empty, error) {
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

	if err := s.ds.User.Activate(user.ID); err != nil {
		s.logger.Errorf("failed to activate user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RPCServer) ResetPassword(ctx context.Context, req *v1.ResetPasswordUserRequest) (*protoempty.Empty, error) {
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

func (s *RPCServer) Get(ctx context.Context, req *protoempty.Empty) (*v1.UserProfile, error) {
	_ = opentracing.SpanFromContext(ctx)

	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	accountProfile, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: user.ID})
	if err != nil {
		s.logger.Errorf("failed to get account profile: %s", err)
	} else {
		userProfile.Account = accountProfile
	}

	return userProfile, nil
}

func (s *RPCServer) GetById(ctx context.Context, req *v1.UserRequest) (*v1.UserProfile, error) {  //nolint
	_ = opentracing.SpanFromContext(ctx)

	user, err := s.ds.User.Get(req.Id)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == datastore.ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.ErrRpcInternal
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	accountProfile, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: user.ID})
	if err != nil {
		s.logger.Errorf("failed to get account profile: %s", err)
	} else {
		userProfile.Account = accountProfile
	}

	return userProfile, nil
}

func (s *RPCServer) Whitelist(ctx context.Context, req *protoempty.Empty) (*v1.WhitelistResponse, error) {
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

func (s *RPCServer) LookupByAddress(ctx context.Context, req *v1.LookupByAddressRequest) (*protoempty.Empty, error) {
	_ = opentracing.SpanFromContext(ctx)

	aReq := &accountsv1.Address{Address: req.Address}
	_, err := s.accounts.GetByAddress(ctx, aReq)
	if err != nil {
		s.logger.Errorf("failed to look up address: %s", err)
		return nil, rpc.ErrRpcNotFound
	}

	return new(protoempty.Empty), nil
}

func (s *RPCServer) Key(ctx context.Context, req *v1.UserRequest) (*v1.KeyResponse, error) {
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

func (s *RPCServer) Keys(ctx context.Context, req *protoempty.Empty) (*v1.KeysResponse, error) {
	_ = opentracing.SpanFromContext(ctx)

	keys, err := s.accounts.Keys(ctx, &protoempty.Empty{})
	if err != nil {
		s.logger.Errorf("failed to get keys: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	items := make([]*v1.KeyResponse, 0)

	for _, i := range keys.Items {
		items = append(items, &v1.KeyResponse{
			Key: i.Key,
		})
	}

	keysResp := &v1.KeysResponse{
		Items: items,
	}

	return keysResp, nil
}

func (s *RPCServer) ListApiTokens(ctx context.Context, req *protoempty.Empty) (*v1.UserApiListResponse, error) {  //nolint
	_ = opentracing.SpanFromContext(ctx)

	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	tokens, err := s.ds.Token.ListByUser(ctx, user.ID)
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

func (s *RPCServer) CreateApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*v1.CreateUserApiTokenResponse, error) {  //nolint
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

	apiToken, err := s.ds.Token.Create(ctx, user.ID, req.Name, token)
	if err != nil {
		s.logger.Errorf("failed to create api token record: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &v1.CreateUserApiTokenResponse{
		Id:    apiToken.ID,
		Name:  apiToken.Name,
		Token: apiToken.Token,
	}, nil
}

func (s *RPCServer) DeleteApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*protoempty.Empty, error) {  //nolint
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

func (s *RPCServer) GetApiToken(ctx context.Context, req *v1.ApiTokenRequest) (*v1.UserApiTokenResponse, error) {  //nolint
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("token", req.Token)

	if req.Token == "" {
		return nil, rpc.ErrRpcNotFound
	}

	token, err := s.ds.Token.GetByToken(ctx, req.Token)
	if err != nil {
		if err == datastore.ErrTokenNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.ErrRpcInternal
	}

	return &v1.UserApiTokenResponse{
		Id:   token.ID,
		Name: token.Name,
	}, nil
}

func (s *RPCServer) StartWithdraw(ctx context.Context, req *v1.StartWithdrawRequest) (*v1.WithdrawResponse, error) {
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

	account, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: user.ID})
	if err != nil {
		return nil, err
	}

	if account.IsLocked {
		s.logger.Error("account is locked")
		return nil, rpc.ErrRpcBadRequest
	}

	famount := new(big.Float)
	famount.SetString(req.Amount)
	if famount == nil || famount.Cmp(big.NewFloat(0)) < 0 {
		s.logger.Errorf("amount is negative or failed to convert: %s", req.Amount)
		return nil, rpc.ErrRpcBadRequest
	}

	amount := new(big.Int).SetBytes([]byte(req.Amount))
	transfer, err := s.accounts.CreateTransfer(ctx, &accountsv1.CreateTransferRequest{
		UserId:    user.ID,
		ToAddress: req.Address,
		Amount:    amount.Bytes(),
	})
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

func (s *RPCServer) Withdraw(ctx context.Context, req *v1.WithdrawRequest) (*protoempty.Empty, error) {
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

	transfer, err := s.accounts.GetTransfer(ctx, &accountsv1.TransferRequest{
		Id: req.TransferId,
	})
	if err != nil {
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

	_, err = s.accounts.ExecuteTransfer(ctx,
		&accountsv1.ExecuteTransferRequest{
			Id:        transfer.Id,
			UserId:    user.ID,
			UserEmail: user.Email,
		},
	)
	if err != nil {
		s.logger.WithError(err).Error("failed to withdraw")
		return nil, rpc.ErrRpcInternal
	}

	return new(protoempty.Empty), nil
}

func (s *RPCServer) createToken(ctx context.Context, user *ds.User, tokenType v1.TokenType) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "createToken")
	defer span.Finish()

	span.SetTag("id", user.ID)
	span.SetTag("email", user.Email)
	span.SetTag("token_type", tokenType)

	claims := auth.ExtendedClaims{
		Type: auth.TokenType(tokenType),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	if tokenType == v1.TokenTypeAPI {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 365).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.authTokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *RPCServer) authenticate(ctx context.Context) (*ds.User, context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, _, err := auth.AuthFromContext(ctx)
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

func (s *RPCServer) getTokenType(ctx context.Context) auth.TokenType {
	tokenType, ok := auth.TypeFromContext(ctx)
	if !ok {
		return auth.TokenType(v1.TokenTypeRegular)
	}

	return tokenType
}
