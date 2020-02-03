package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr               string `default:"0.0.0.0:5000" envconfig:"RPC_ADDR"`
	AccountsRPCAddr       string `default:"0.0.0.0:5001" envconfig:"ACCOUNTS_RPC_ADDR"`
	ServiceManagerRPCAddr string `default:"0.0.0.0:5017" envconfig:"SERVICE_MANAGER_RPC_ADDR"`
	DBURI                 string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local" envconfig:"DBURI"`
	MQURI                 string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`
	AuthTokenSecret       string `default:"secret" envconfig:"AUTH_TOKEN_SECRET"`
	AuthRecoverySecret    string `default:"secret" envconfig:"AUTH_RECOVERY_SECRET"`

	Logger *logrus.Entry `envconfig:"-"`
}
