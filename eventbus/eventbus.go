package eventbus

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/videocoin/cloud-pkg/mqmux"

	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	notificationv1 "github.com/videocoin/cloud-api/notifications/v1"
)

type EventBus struct {
	mq     *mqmux.WorkerMux
	logger *logrus.Entry
}

type Config struct {
	Logger *logrus.Entry
	URI    string
	Name   string
}

func New(c *Config) (*EventBus, error) {
	mq, err := mqmux.NewWorkerMux(c.URI, c.Name)
	if err != nil {
		return nil, err
	}

	return &EventBus{
		logger: c.Logger,
		mq:     mq,
	}, nil
}

func (e *EventBus) Start() error {
	if err := e.registerPublishers(); err != nil {
		return err
	}

	if err := e.registerConsumers(); err != nil {
		return err
	}

	return e.mq.Run()
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) registerPublishers() error {
	if err := e.mq.Publisher("accounts.create"); err != nil {
		return err
	}

	if err := e.mq.Publisher("notifications.send"); err != nil {
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

	if err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	); err != nil {
		return err
	}

	return e.mq.PublishX("accounts.create", req, headers)
}

func (e *EventBus) SendNotification(span opentracing.Span, req *notificationv1.Notification) error {
	headers := make(amqp.Table)
	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "users")

	if err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		mqmux.RMQHeaderCarrier(headers),
	); err != nil {
		return err
	}

	return e.mq.PublishX("notifications.send", req, headers)
}
