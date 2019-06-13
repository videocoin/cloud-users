package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	notificationv1 "github.com/videocoin/cloud-api/notifications/v1"
	"github.com/videocoin/cloud-pkg/mqmux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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

func (e *EventBus) CreateUserAccount(span opentracing.Span, req *accountsv1.AccountRequest) error {
	headers := make(amqp.Table)
	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "users")

	span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	)

	return e.mq.PublishX("account/create", req, headers)
}

func (e *EventBus) SendNotification(span opentracing.Span, req *notificationv1.Notification) error {
	headers := make(amqp.Table)
	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "users")

	span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	)

	return e.mq.PublishX("notifications/send", req, headers)
}
