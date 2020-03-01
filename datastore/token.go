package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/opentracing/opentracing-go"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-pkg/uuid4"
)

var (
	ErrTokenNotFound = errors.New("token is not found")
)

type TokenDatastore struct {
	db *gorm.DB
}

type UserAPIToken struct {
	ID        string     `gorm:"type:varchar(36);PRIMARY_KEY"`
	UserID    string     `gorm:"type:varchar(36);DEFAULT:null"`
	Name      string     `gorm:"type:varchar(100);DEFAULT:null"`
	Token     string     `gorm:"type:varchar(255);DEFAULT:null"`
	CreatedAt *time.Time `gorm:"type:timestamp NULL;DEFAULT:null"`
}

func NewTokenDatastore(db *gorm.DB) (*TokenDatastore, error) {
	return &TokenDatastore{db: db}, nil
}

func (ds *TokenDatastore) Create(ctx context.Context, userID, name, token string) (*UserAPIToken, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	span.SetTag("user_id", userID)
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

	apiToken := &UserAPIToken{
		ID:        id,
		UserID:    userID,
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

func (ds *TokenDatastore) ListByUser(ctx context.Context, userID string) ([]*UserAPIToken, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListByUser")
	defer span.Finish()

	span.SetTag("user_id", userID)

	tokens := []*UserAPIToken{}
	if err := ds.db.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to get user api tokens: %s", err)
	}

	return tokens, nil
}

func (ds *TokenDatastore) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("id", id)

	token := &UserAPIToken{
		ID: id,
	}
	if err := ds.db.Delete(token).Error; err != nil {
		return fmt.Errorf("failed to delete user api token: %s", err)
	}

	return nil
}

func (ds *TokenDatastore) GetByToken(ctx context.Context, token string) (*UserAPIToken, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetByToken")
	defer span.Finish()

	span.SetTag("token", token)

	t := &UserAPIToken{}
	if err := ds.db.Where("token = ?", token).First(t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTokenNotFound
		}

		return nil, fmt.Errorf("failed to get api token by token: %s", err)
	}

	return t, nil
}
