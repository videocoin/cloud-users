package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/jinzhu/gorm"
	v1 "github.com/videocoin/cloud-api/transfers/v1"
	"github.com/videocoin/cloud-pkg/uuid4"
)

type TransferDatastore struct {
	db *gorm.DB
}

func NewTransferDatastore(db *gorm.DB) (*TransferDatastore, error) {
	db.AutoMigrate(&v1.Transfer{})
	return &TransferDatastore{db: db}, nil
}

func (ds *TransferDatastore) Create(ctx context.Context, userId, address string, amount float64) (*v1.Transfer, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	tx := ds.db.Begin()

	id, err := uuid4.New()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	createdAt, err := ptypes.Timestamp(ptypes.TimestampNow())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	ts, err := ptypes.TimestampProto(time.Now().Add(time.Minute * 10))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	expiresAt, err := ptypes.Timestamp(ts)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transfer := &v1.Transfer{
		Id:        id,
		UserId:    userId,
		Kind:      v1.TransferKindWithdraw,
		Pin:       encodeToString(6),
		ToAddress: address,
		Amount:    amount,
		CreatedAt: &createdAt,
		ExpiresAt: &expiresAt,
		// ? amount ?
	}

	if err = tx.Create(transfer).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return transfer, nil
}

func (ds *TransferDatastore) Get(ctx context.Context, id string) (*v1.Transfer, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("id", id)

	transfer := new(v1.Transfer)
	if err := ds.db.Where("id = ?", id).First(&transfer).Error; err != nil {
		return nil, fmt.Errorf("failed to get transfer: %s", err)
	}

	return transfer, nil
}

func (ds *TransferDatastore) ListByUser(ctx context.Context, userId string) ([]*v1.Transfer, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListByUser")
	defer span.Finish()

	span.SetTag("user_id", userId)

	transfers := []*v1.Transfer{}
	if err := ds.db.Where("user_id = ?", userId).Find(&transfers).Error; err != nil {
		return nil, fmt.Errorf("failed to get user transfers: %s", err)
	}

	return transfers, nil
}

func (ds *TransferDatastore) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("id", id)

	transfer := &v1.Transfer{
		Id: id,
	}
	if err := ds.db.Delete(transfer).Error; err != nil {
		return fmt.Errorf("failed to delete user transfer: %s", err)
	}

	return nil
}
