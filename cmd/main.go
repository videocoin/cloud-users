package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"github.com/kelseyhightower/envconfig"
	pkglogger "github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-users/service"
)

var (
	ServiceName string = "users"
	Version     string = "dev"
)

func main() {
	logger := pkglogger.NewLogrusLogger(ServiceName, Version, nil)

	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		logger.Info(err.Error())
	} else {
		defer closer.Close()
	}

	cfg := &service.Config{
		Name:    ServiceName,
		Version: Version,
	}

	err = envconfig.Process(ServiceName, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := ctxlogrus.ToContext(context.Background(), logger)
	svc, err := service.NewService(ctx, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	errCh := make(chan error, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		logger.Infof("received signal %s", sig)
		exit <- true
	}()

	logger.Info("starting")
	go svc.Start(errCh)

	select {
	case <-exit:
		break
	case err := <-errCh:
		if err != nil {
			logger.Error(err)
		}
		break
	}

	logger.Info("stopping")
	err = svc.Stop()
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info("stopped")
}
