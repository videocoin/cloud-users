package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr         string `default:"0.0.0.0:5000"`
	AccountsRPCAddr string `default:"0.0.0.0:5001"`
	DBURI           string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local" envconfig:"DBURI"`
	MQURI           string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`
	Secret          string `default:"secret" envconfig:"SECRET"`
	RecoverySecret  string `default:"secret" envconfig:"RECOVERYSECRET"`
	CentSecret      string `default:"secret" envconfig:"CENTSECRET"`

	Logger *logrus.Entry `envconfig:"-"`
}
