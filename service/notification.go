package service

import (
	"context"

	notificationv1 "github.com/VideoCoin/cloud-api/notifications/v1"
	v1 "github.com/VideoCoin/cloud-api/users/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
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

func (c *NotificationClient) SendEmailWelcome(ctx context.Context, user *v1.User) error {
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

	err := c.eb.SendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotificationClient) SendEmailRecovery(ctx context.Context, user *v1.User, token string) error {
	params := map[string]string{
		"to":     user.Email,
		"domain": "videocoin.network",
		"token":  token,
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_recovery",
		Params:   params,
	}

	err := c.eb.SendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotificationClient) SendTestPush(ctx context.Context, user *v1.User) error {
	params := map[string]string{
		"event":   "user/created",
		"user_id": user.Id,
		"name":    user.Name,
	}

	notification := &notificationv1.Notification{
		Target: notificationv1.NotificationTarget_WEB,
		Params: params,
	}

	err := c.eb.SendNotification(notification)
	if err != nil {
		return err
	}

	return nil
}
