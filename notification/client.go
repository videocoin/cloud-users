package notification

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/notifications/v1"
	ds "github.com/videocoin/cloud-users/datastore"
	"github.com/videocoin/cloud-users/eventbus"
)

type Client struct {
	eb     *eventbus.EventBus
	logger *logrus.Entry
}

func NewClient(eb *eventbus.EventBus, logger *logrus.Entry) (*Client, error) {
	return &Client{
		eb:     eb,
		logger: logger,
	}, nil
}

func (c *Client) SendEmailWaitlisted(ctx context.Context, user *ds.User) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailWaitlisted")
	defer span.Finish()

	params := map[string]string{
		"to": user.Email,
	}

	notification := &v1.Notification{
		Target:   v1.NotificationTarget_EMAIL,
		Template: "user_waitlisted",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendEmailWelcome(ctx context.Context, user *ds.User) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailWelcome")
	defer span.Finish()

	params := map[string]string{
		"to":   user.Email,
		"name": user.FirstName + " " + user.LastName,
	}

	notification := &v1.Notification{
		Target:   v1.NotificationTarget_EMAIL,
		Template: "user_welcome",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendEmailRecovery(ctx context.Context, user *ds.User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailRecovery")
	defer span.Finish()

	params := map[string]string{
		"to":    user.Email,
		"token": token,
	}

	notification := &v1.Notification{
		Target:   v1.NotificationTarget_EMAIL,
		Template: "user_recovery",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendEmailConfirmation(ctx context.Context, user *ds.User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailConfirmation")
	defer span.Finish()

	params := map[string]string{
		"to":    user.Email,
		"token": token,
	}

	notification := &v1.Notification{
		Target:   v1.NotificationTarget_EMAIL,
		Template: "user_confirmation",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}
