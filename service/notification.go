package service

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	notificationv1 "github.com/videocoin/cloud-api/notifications/v1"
	v1 "github.com/videocoin/cloud-api/users/v1"
)

type NotificationClient struct {
	eb     *EventBus
	logger *logrus.Entry
}

func NewNotificationClient(eb *EventBus, logger *logrus.Entry) (*NotificationClient, error) {
	return &NotificationClient{
		eb:     eb,
		logger: logger,
	}, nil
}

func (c *NotificationClient) SendEmailWaitlisted(ctx context.Context, user *v1.User) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailWaitlisted")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	params := map[string]string{
		"to":     user.Email,
		"domain": md.Get("x-forwarded-host"),
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_waitlisted",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotificationClient) SendEmailWelcome(ctx context.Context, user *v1.User) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	params := map[string]string{
		"to":     user.Email,
		"name":   user.Name,
		"domain": md.Get("x-forwarded-host"),
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_welcome",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotificationClient) SendEmailRecovery(ctx context.Context, user *v1.User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	params := map[string]string{
		"to":     user.Email,
		"token":  token,
		"domain": md.Get("x-forwarded-host"),
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_recovery",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}
