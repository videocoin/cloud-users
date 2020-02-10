package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	smv1 "github.com/videocoin/cloud-api/servicemanager/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/mqmux"
	"github.com/videocoin/cloud-users/datastore"
)

type Service struct {
	cfg *Config
	rpc *RpcServer
	eb  *EventBus
}

func NewService(cfg *Config) (*Service, error) {
	ds, err := datastore.NewDatastore(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	conn, err := grpcutil.Connect(cfg.AccountsRPCAddr, cfg.Logger.WithField("system", "accountcli"))
	if err != nil {
		return nil, err
	}
	accounts := accountsv1.NewAccountServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.ServiceManagerRPCAddr, cfg.Logger.WithField("system", "smcli"))
	if err != nil {
		return nil, err
	}
	sm := smv1.NewServiceManagerClient(conn)

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

	rpcConfig := &RpcServerOptions{
		Addr:               cfg.RPCAddr,
		AuthTokenSecret:    cfg.AuthTokenSecret,
		AuthRecoverySecret: cfg.AuthRecoverySecret,
		Logger:             cfg.Logger,
		DS:                 ds,
		EB:                 eb,
		Accounts:           accounts,
		Sm:                 sm,
	}

	rpc, err := NewRpcServer(rpcConfig)
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
	go s.rpc.Start()
	go s.eb.Start()
	return nil
}

func (s *Service) Stop() error {
	s.eb.Stop()
	return nil
}
