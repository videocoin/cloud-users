package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr         string `default:"127.0.0.1:5000"`
	AccountsRPCAddr string `default:"127.0.0.1:2000" envconfig:"VC_ACCOUNTS_RPCADDR"`
	AccountsEnabled bool   `default:"false" envconfig:"VC_ACCOUNT_ENABLE"`
	DBURI           string `default:"mysql:mysql@/vc-user?charset=utf8&parseTime=True&loc=Local"`
	MQURI           string `default:"amqp://rabbitmq:bitnami@127.0.0.1:5672" envconfig:"VC_MQURI"`
	Secret          string `default:"secret"`

	Logger *logrus.Entry `envconfig:"-"`
}
