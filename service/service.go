package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/mqmux"
	"github.com/videocoin/cloud-users/datastore"
	"google.golang.org/grpc"
)

type Service struct {
	cfg *Config
	rpc *RPCServer
	eb  *EventBus
}

func NewService(cfg *Config) (*Service, error) {
	ds, err := datastore.NewDatastore(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	alogger := cfg.Logger.WithField("system", "accountcli")
	aGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(alogger)
	accountsConn, err := grpc.Dial(cfg.AccountsRPCAddr, aGrpcDialOpts...)
	if err != nil {
		return nil, err
	}

	accounts := accountsv1.NewAccountServiceClient(accountsConn)

	mq, err := mqmux.NewWorkerMux(cfg.MQURI, cfg.Name)
	if err != nil {
		return nil, err
	}
	mq.Logger = cfg.Logger.WithField("system", "mq")

	eblogger := cfg.Logger.WithField("system", "eventbus")
	eb, err := NewEventBus(mq, eblogger)
	if err != nil {
		return nil, err
	}

	rpcConfig := &RPCServerOptions{
		Addr:               cfg.RPCAddr,
		AuthTokenSecret:    cfg.AuthTokenSecret,
		AuthRecoverySecret: cfg.AuthRecoverySecret,
		Logger:             cfg.Logger,
		DS:                 ds,
		EB:                 eb,
		Accounts:           accounts,
	}

	rpc, err := NewRPCServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg: cfg,
		rpc: rpc,
		eb:  eb,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()  //nolint
	go s.eb.Start()  //nolint
	return nil
}

func (s *Service) Stop() error {
	err := s.eb.Stop()
	return err
}
