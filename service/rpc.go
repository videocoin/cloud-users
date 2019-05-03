package service

import (
	"context"
	"fmt"
	"net"
	"time"

	accountsv1 "github.com/VideoCoin/cloud-api/accounts/v1"
	"github.com/VideoCoin/cloud-api/rpc"
	v1 "github.com/VideoCoin/cloud-api/users/v1"
	"github.com/VideoCoin/cloud-pkg/auth"
	"github.com/VideoCoin/cloud-pkg/grpcutil"
	jwt "github.com/dgrijalva/jwt-go"
	protoempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcServerOptions struct {
	Addr           string
	Secret         string
	RecoverySecret string
	Logger         *logrus.Entry
	DS             *Datastore
	Accounts       accountsv1.AccountServiceClient
	EB             *EventBus
}

type RpcServer struct {
	addr           string
	secret         string
	recoverySecret string
	grpc           *grpc.Server
	listen         net.Listener
	logger         *logrus.Entry
	ds             *Datastore
	eb             *EventBus
	accounts       accountsv1.AccountServiceClient
	notifications  *NotificationClient
}

func NewRpcServer(opts *RpcServerOptions) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)

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
		addr:          opts.Addr,
		secret:        opts.Secret,
		grpc:          grpcServer,
		listen:        listen,
		logger:        opts.Logger,
		ds:            opts.DS,
		eb:            opts.EB,
		accounts:      opts.Accounts,
		notifications: nc,
	}

	v1.RegisterUserServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) Health(ctx context.Context, req *protoempty.Empty) (*rpc.HealthStatus, error) {
	return &rpc.HealthStatus{Status: "OK"}, nil
}

func (s *RpcServer) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.LoginUserResponse, error) {
	resp := new(v1.LoginUserResponse)

	err := req.Validate()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if req.Password != req.ConfirmPassword {
		respErr := &rpc.MultiValidationError{
			Errors: []*rpc.ValidationError{
				&rpc.ValidationError{
					Field:   "password",
					Message: "Passwords do not match",
				},
			},
		}
		return resp, rpc.NewRpcValidationError(respErr)
	}

	user, err := s.ds.User.Register(req.Email, req.Name, req.Password)
	if err != nil {
		if err == ErrUserIsAlreadyExists {
			respErr := &rpc.MultiValidationError{
				Errors: []*rpc.ValidationError{
					&rpc.ValidationError{
						Field:   "email",
						Message: "User is already registered",
					},
				},
			}
			return resp, rpc.NewRpcValidationError(respErr)
		}

		s.logger.Error(err)

		return nil, rpc.ErrRpcInternal
	}

	token, err := s.createToken(ctx, user)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.ds.User.UpdateAuthToken(user, token)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	resp.Token = user.Token

	accReq := &accountsv1.AccountRequest{OwnerID: user.Id}
	_, err = s.accounts.Create(ctx, accReq)
	if err != nil {
		s.logger.Warningf("failed to create account: %s", err)
		accErr := s.eb.CreateUserAccount(accReq)
		if accErr != nil {
			s.logger.Errorf("failed to create account via eventbus: %s", accErr)
		}
	}

	err = s.notifications.SendEmailWelcome(ctx, user)
	if err != nil {
		s.logger.WithField("faield to send welcome email to user id", user.Id).Error(err)
	}

	err = s.notifications.SendTestPush(ctx, user)
	if err != nil {
		s.logger.WithField("faield to send push to user id", user.Id).Error(err)
	}

	return resp, nil
}

func (s *RpcServer) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.LoginUserResponse, error) {
	err := req.Validate()
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	user, err := s.ds.User.GetByEmail(req.Email)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		return nil, rpc.ErrRpcInternal
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return nil, rpc.ErrRpcUnauthenticated
	}

	token, err := s.createToken(ctx, user)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.ds.User.UpdateAuthToken(user, token)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	resp := &v1.LoginUserResponse{
		Token: user.Token,
	}

	return resp, nil
}

func (s *RpcServer) Logout(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
	user, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.ds.User.ResetAuthToken(user)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) StartRecovery(ctx context.Context, req *v1.StartRecoveryUserRequest) (*protoempty.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.ds.User.GetByEmail(req.Email)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		return nil, rpc.ErrRpcInternal
	}

	token := newRecoveryToken(req.Email, 12*time.Hour, []byte(user.Password), []byte(s.recoverySecret))

	err = s.notifications.SendEmailRecovery(ctx, user, token)
	if err != nil {
		s.logger.WithField("faield to send recovery email to user id", user.Id).Error(err)
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Recover(ctx context.Context, req *v1.RecoverUserRequest) (*protoempty.Empty, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	user, err := verifyRecoveryToken(req.Token, s.ds.User.GetByEmail, []byte(s.recoverySecret))
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		if err == ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		return nil, rpc.ErrRpcInternal
	}

	err = s.ds.User.ResetPassword(user, req.Password)
	if err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) ResetPassword(ctx context.Context, req *v1.ResetPasswordUserRequest) (*protoempty.Empty, error) {
	fmt.Printf("received request 2: %v", req)
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.ds.User.ResetPassword(user, req.Password)
	if err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Get(ctx context.Context, req *protoempty.Empty) (*v1.UserProfile, error) {
	user, ctx, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	userProfile := new(v1.UserProfile)
	err = copier.Copy(userProfile, user)
	if err != nil {
		return nil, rpc.ErrRpcInternal
	}

	accountProfile, err := s.accounts.Get(ctx, &accountsv1.AccountRequest{OwnerID: user.Id})
	if err != nil {
		s.logger.Errorf("failed to get account profile: %s", err)
	} else {
		userProfile.Account = accountProfile
	}

	return userProfile, nil
}

func (s *RpcServer) Whitelist(ctx context.Context, req *protoempty.Empty) (*v1.WhitelistResponse, error) {
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
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	aReq := &accountsv1.Address{Address: req.Address}
	_, err = s.accounts.GetByAddress(ctx, aReq)
	if err != nil {
		s.logger.Errorf("failed to look up address: %s", err)
		return nil, rpc.ErrRpcNotFound

	}

	return new(protoempty.Empty), nil
}

func (s *RpcServer) createToken(ctx context.Context, user *v1.User) (string, error) {
	claims := &auth.JWTClaims{
		UserID: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *RpcServer) authenticate(ctx context.Context) (*v1.User, context.Context, error) {
	ctx = auth.NewContextWithSecretKey(ctx, s.secret)
	ctx, err := auth.AuthFromContext(ctx)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
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
