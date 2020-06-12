package rpc

import (
	"net"

	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	v1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-users/datastore"
	"github.com/videocoin/cloud-users/eventbus"
	"github.com/videocoin/cloud-users/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOptions struct {
	Addr               string
	AuthTokenSecret    string
	AuthRecoverySecret string

	Logger   *logrus.Entry
	DS       *datastore.Datastore
	Accounts accountsv1.AccountServiceClient
	EB       *eventbus.EventBus
}

type Server struct {
	addr               string
	authTokenSecret    string
	authRecoverySecret string

	server        *grpc.Server
	listen        net.Listener
	logger        *logrus.Entry
	ds            *datastore.Datastore
	eb            *eventbus.EventBus
	accounts      accountsv1.AccountServiceClient
	notifications *notification.Client
	validator     *requestValidator
}

func NewServer(opts *ServerOptions) (*Server, error) {
	server := grpc.NewServer(grpcutil.DefaultServerOpts(opts.Logger)...)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	nc, err := notification.NewClient(
		opts.EB, opts.Logger.WithField("system", "notification"))
	if err != nil {
		return nil, err
	}

	validator, err := newRequestValidator()
	if err != nil {
		return nil, err
	}

	self := &Server{
		addr:               opts.Addr,
		authTokenSecret:    opts.AuthTokenSecret,
		authRecoverySecret: opts.AuthRecoverySecret,
		server:             server,
		listen:             listen,
		logger:             opts.Logger,
		ds:                 opts.DS,
		eb:                 opts.EB,
		accounts:           opts.Accounts,
		notifications:      nc,
		validator:          validator,
	}

	v1.RegisterUserServiceServer(server, self)
	reflection.Register(server)

	return self, nil
}

func (s *Server) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.server.Serve(s.listen)
}
