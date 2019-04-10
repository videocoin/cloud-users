package service

import (
	"context"
	"net"
	"time"

	accountsv1 "github.com/VideoCoin/cloud-api/accounts/v1"
	"github.com/VideoCoin/cloud-api/rpc"
	"github.com/VideoCoin/cloud-api/users/v1"
	"github.com/VideoCoin/cloud-api/validator"
	"github.com/VideoCoin/cloud-pkg/grpcutil"
	"github.com/VideoCoin/cloud-users/auth"
	jwt "github.com/dgrijalva/jwt-go"
	protoempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcServerOptions struct {
	Addr            string
	Secret          string
	Logger          *logrus.Entry
	DS              *Datastore
	Accounts        accountsv1.AccountServiceClient
	AccountsEnabled bool
	EB              *EventBus
}

type RpcServer struct {
	addr            string
	secret          string
	grpc            *grpc.Server
	listen          net.Listener
	logger          *logrus.Entry
	ds              *Datastore
	eb              *EventBus
	accounts        accountsv1.AccountServiceClient
	accountsEnabled bool
}

func NewRpcServer(opts *RpcServerOptions) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		addr:     opts.Addr,
		secret:   opts.Secret,
		grpc:     grpcServer,
		listen:   listen,
		logger:   opts.Logger,
		ds:       opts.DS,
		eb:       opts.EB,
		accounts: opts.Accounts,
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

	verr := validator.ValidateCreateUserRequest(req)
	if verr != nil {
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.Register(req.Email, req.Password)
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

	if s.accountsEnabled {
		accReq := &accountsv1.AccountRequest{OwnerID: user.Id}
		_, err := s.accounts.Create(ctx, accReq)
		if err != nil {
			s.logger.Warningf("failed to create user: %s", err)
			billErr := s.eb.CreateUserAccount(accReq)
			if billErr != nil {
				s.logger.Errorf("failed to create user via eventbus: %s", billErr)
			}
		}
	}

	return resp, nil
}

func (s *RpcServer) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.LoginUserResponse, error) {
	verr := validator.ValidateLoginUserRequest(req)
	if verr != nil {
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.GetByEmail(req.Email)
	if err != nil {
		return nil, rpc.ErrRpcUnauthenticated
	}

	if !CheckPasswordHash(req.Password, user.Password) {
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
	empty := &protoempty.Empty{}

	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = s.ds.User.ResetAuthToken(user)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return empty, nil
}

func (s *RpcServer) GetUserProfile(ctx context.Context, req *protoempty.Empty) (*v1.UserProfile, error) {
	user, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	userProfile := new(v1.UserProfile)
	copier.Copy(userProfile, user)

	if s.accountsEnabled {
		accountProfile, err := s.accounts.Get(ctx, &accountsv1.AccountRequest{OwnerID: user.Id})
		if err != nil {
			s.logger.Errorf("failed to get account profile: %s", err)
		} else {
			userProfile.Account = accountProfile
		}
	}

	return userProfile, nil
}

func (s *RpcServer) GetUserById(ctx context.Context, req *v1.UserRequest) (*v1.UserProfile, error) {
	user, err := s.ds.User.GetByID(req.Id)
	if err != nil {
		s.logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	userProfile := new(v1.UserProfile)
	copier.Copy(userProfile, user)

	if s.accountsEnabled {
		accountProfile, err := s.accounts.Get(ctx, &accountsv1.AccountRequest{OwnerID: user.Id})
		if err != nil {
			s.logger.Errorf("failed to get account profile: %s", err)
		} else {
			userProfile.Account = accountProfile
		}
	}

	return userProfile, nil
}

func (s *RpcServer) GetList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	resp := &v1.ListResponse{
		Items: []*v1.User{},
	}

	users, err := s.ds.User.GetList()
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	resp.Items = users

	return resp, nil
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

	user, err := s.ds.User.GetByID(userID)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	if user.Token == "" {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	return user, ctx, nil
}
