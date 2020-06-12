package service

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-users/datastore"
	"github.com/videocoin/cloud-users/eventbus"
	"github.com/videocoin/cloud-users/rpc"
)

type Service struct {
	cfg    *Config
	server *rpc.Server
	eb     *eventbus.EventBus
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	ds, err := datastore.New(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	sc, err := clientv1.NewServiceClientFromEnvconfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	eb, err := eventbus.New(&eventbus.Config{
		URI:    cfg.MQURI,
		Name:   cfg.Name,
		Logger: ctxlogrus.Extract(ctx).WithField("system", "eventbus"),
	})
	if err != nil {
		return nil, err
	}

	rpcConfig := &rpc.ServerOptions{
		Addr:               cfg.RPCAddr,
		AuthTokenSecret:    cfg.AuthTokenSecret,
		AuthRecoverySecret: cfg.AuthRecoverySecret,
		Logger:             cfg.Logger,
		DS:                 ds,
		EB:                 eb,
		Accounts:           sc.Accounts,
	}

	server, err := rpc.NewServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:    cfg,
		server: server,
		eb:     eb,
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.server.Start()
	}()

	go func() {
		errCh <- s.eb.Start()
	}()
}

func (s *Service) Stop() error {
	err := s.eb.Stop()
	return err
}
