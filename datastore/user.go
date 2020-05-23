package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/opentracing/opentracing-go"

	"github.com/jinzhu/gorm"
	v1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/dbutil"
	"github.com/videocoin/cloud-pkg/uuid4"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserDatastore struct {
	db *gorm.DB
}

type User struct {
	ID          string
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Country     string
	Region      string
	City        string
	Zip         string
	Address1    string `gorm:"Column:address_1"`
	Address2    string `gorm:"Column:address_2"`
	Role        v1.UserRole
	UIRole      v1.UserUIRole `gorm:"Column:uirole"`
	IsActive    bool          `gorm:"Column:is_active"`
	ActivatedAt *time.Time
	CreatedAt   *time.Time
	Token       string
}

func NewUserDatastore(db *gorm.DB) (*UserDatastore, error) {
	return &UserDatastore{db: db}, nil
}

func (ds *UserDatastore) List() ([]*User, error) {
	users := []*User{}

	if err := ds.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users list: %s", err)
	}

	return users, nil
}

func (ds *UserDatastore) Get(id string) (*User, error) {
	user := &User{}

	if err := ds.db.Where("id = ?", id).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) GetByEmail(ctx context.Context, email string) (*User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetByEmail")
	defer span.Finish()

	span.SetTag("email", email)

	user := &User{}
	if err := ds.db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) GetByVerificationCode(code string) (*User, error) {
	user := &User{}
	if err := ds.db.Where("verification_code = ?", code).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by verification code: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) Validate(ctx context.Context, email string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Validate")
	defer span.Finish()

	span.SetTag("email", email)

	tx := ds.db.Begin()

	user := &User{}
	err := tx.Where("email = ?", email).First(user).Error
	if err == nil {
		tx.Rollback()
		return ErrUserAlreadyExists
	}

	return nil
}

func (ds *UserDatastore) Register(ctx context.Context, req *v1.CreateUserRequest) (*User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Register")
	defer span.Finish()

	span.SetTag("email", req.Email)

	tx := ds.db.Begin()

	user := &User{}
	err := tx.Where("email = ?", req.Email).First(user).Error
	if err == nil {
		tx.Rollback()
		return nil, ErrUserAlreadyExists
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	id, err := uuid4.New()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	passwordHash, _ := hashPassword(ctx, req.Password)
	time, err := ptypes.Timestamp(ptypes.TimestampNow())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user = &User{
		ID:        id,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Country:   req.Country,
		Region:    req.Region,
		City:      req.City,
		Zip:       req.Zip,
		Address1:  req.Address_1,
		Address2:  req.Address_2,
		Password:  passwordHash,
		IsActive:  false,
		CreatedAt: &time,
		Role:      v1.UserRoleMiner,
		UIRole:    req.UiRole,
	}

	if err = tx.Create(user).Error; err != nil {
		ec := dbutil.ErrorCode(err)
		if ec == dbutil.ErrDuplicateEntry {
			return nil, ErrUserAlreadyExists
		}

		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return user, nil
}

func (ds *UserDatastore) ResetPassword(ctx context.Context, user *User, password string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Recover")
	defer span.Finish()

	span.SetTag("email", user.Email)

	passwordHash, _ := hashPassword(ctx, password)
	user.Password = passwordHash

	updates := map[string]interface{}{
		"password": user.Password,
	}

	if err := ds.db.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (ds *UserDatastore) UpdateAuthToken(ctx context.Context, user *User, token string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UpdateAuthToken")
	defer span.Finish()

	user.Token = token
	updates := map[string]interface{}{
		"token": user.Token,
	}

	if err := ds.db.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (ds *UserDatastore) ResetAuthToken(user *User) error {
	user.Token = ""

	updates := map[string]interface{}{
		"token": user.Token,
	}

	if err := ds.db.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (ds *UserDatastore) Activate(userID string) error {
	time, err := ptypes.Timestamp(ptypes.TimestampNow())
	if err != nil {
		return err
	}

	user := &User{
		ID: userID,
	}

	updates := map[string]interface{}{
		"is_active":   true,
		"activatedAt": &time,
	}

	if err = ds.db.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (ds *UserDatastore) UpdateUIRole(ctx context.Context, user *User, uiRole v1.UserUIRole) error {
	updates := map[string]interface{}{
		"ui_role": uiRole,
	}

	if err := ds.db.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
