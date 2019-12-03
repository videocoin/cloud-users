package service

import (
	"context"
	"math/big"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	notificationv1 "github.com/videocoin/cloud-api/notifications/v1"
	"github.com/videocoin/cloud-pkg/ethutils"
	ds "github.com/videocoin/cloud-users/datastore"
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

func (c *NotificationClient) SendEmailWaitlisted(ctx context.Context, user *ds.User) error {
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

func (c *NotificationClient) SendEmailWelcome(ctx context.Context, user *ds.User) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailWelcome")
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

func (c *NotificationClient) SendEmailRecovery(ctx context.Context, user *ds.User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailRecovery")
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

func (c *NotificationClient) SendEmailConfirmation(ctx context.Context, user *v1.User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendEmailConfirmation")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	params := map[string]string{
		"to":     user.Email,
		"token":  token,
		"domain": md.Get("x-forwarded-host"),
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_confirmation",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotificationClient) SendWithdrawTransfer(ctx context.Context, user *ds.User, transfer *accountsv1.TransferResponse) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SendWithdrawTransfer")
	defer span.Finish()

	md := metautils.ExtractIncoming(ctx)

	amount := new(big.Int)
	amount, _ = amount.SetString(string(transfer.Amount), 10)
	vdc, _ := ethutils.WeiToEth(amount)

	params := map[string]string{
		"to":      user.Email,
		"address": transfer.ToAddress,
		"amount":  vdc.String(),
		"pin":     transfer.Pin,
		"domain":  md.Get("x-forwarded-host"),
	}

	notification := &notificationv1.Notification{
		Target:   notificationv1.NotificationTarget_EMAIL,
		Template: "user_withdraw_confirmation",
		Params:   params,
	}

	err := c.eb.SendNotification(span, notification)
	if err != nil {
		return err
	}

	return nil
}
