package service

import (
	accountsv1 "github.com/VideoCoin/cloud-api/accounts/v1"
	notificationv1 "github.com/VideoCoin/cloud-api/notifications/v1"
	"github.com/VideoCoin/cloud-pkg/mqmux"
	"github.com/sirupsen/logrus"
)

type EventBus struct {
	mq     *mqmux.WorkerMux
	logger *logrus.Entry
}

func NewEventBus(mq *mqmux.WorkerMux, logger *logrus.Entry) (*EventBus, error) {
	return &EventBus{
		logger: logger,
		mq:     mq,
	}, nil
}

func (e *EventBus) Start() error {
	err := e.registerPublishers()
	if err != nil {
		return err
	}

	err = e.registerConsumers()
	if err != nil {
		return err
	}

	return e.mq.Run()
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) registerPublishers() error {
	if err := e.mq.Publisher("account/create"); err != nil {
		return err
	}

	if err := e.mq.Publisher("notifications/send"); err != nil {
		return err
	}

	if err := e.mq.Publisher("notifications/send"); err != nil {
		return err
	}

	return nil
}

func (e *EventBus) registerConsumers() error {
	return nil
}

func (e *EventBus) CreateUserAccount(req *accountsv1.AccountRequest) error {
	return e.mq.Publish("account/create", req)
}

func (e *EventBus) SendNotification(req *notificationv1.Notification) error {
	return e.mq.Publish("notifications/send", req)
}
