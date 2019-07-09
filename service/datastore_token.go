package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/opentracing/opentracing-go"

	"github.com/jinzhu/gorm"
	v1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/uuid4"
)

type TokenDatastore struct {
	db *gorm.DB
}

func NewTokenDatastore(db *gorm.DB) (*TokenDatastore, error) {
	db.AutoMigrate(&v1.UserApiToken{})
	return &TokenDatastore{db: db}, nil
}

func (ds *TokenDatastore) Create(ctx context.Context, userId, name, token string) (*v1.UserApiToken, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	span.SetTag("user_id", userId)
	span.SetTag("name", name)

	tx := ds.db.Begin()

	id, err := uuid4.New()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	time, err := ptypes.Timestamp(ptypes.TimestampNow())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	apiToken := &v1.UserApiToken{
		Id:        id,
		UserId:    userId,
		Name:      name,
		Token:     token,
		CreatedAt: &time,
	}

	if err = tx.Create(apiToken).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return apiToken, nil
}

func (ds *TokenDatastore) ListByUser(ctx context.Context, userId string) ([]*v1.UserApiToken, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("user_id", userId)

	tokens := []*v1.UserApiToken{}

	if err := ds.db.Where("user_id = ?", userId).Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to get user api tokens: %s", err)
	}

	return tokens, nil
}

func (ds *TokenDatastore) Delete(ctx context.Context, tokenId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("token_id", tokenId)

	apiToken := &v1.UserApiToken{
		Id: tokenId,
	}
	if err := ds.db.Delete(apiToken).Error; err != nil {
		return fmt.Errorf("failed to delete user api token: %s", err)
	}

	return nil
}
