package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/VideoCoin/cloud-api/users/v1"
	"github.com/VideoCoin/cloud-users/pkg/dbutil"
	"github.com/VideoCoin/cloud-users/pkg/uuid4"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserIsAlreadyExists = errors.New("user is already exists")
)

type UserDatastore struct {
	db *gorm.DB
}

func NewUserDatastore(db *gorm.DB) (*UserDatastore, error) {
	db.AutoMigrate(&v1.User{})
	return &UserDatastore{db: db}, nil
}

func (ds *UserDatastore) GetList() ([]*v1.User, error) {
	users := []*v1.User{}

	err := ds.db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users list: %s", err)
	}

	return users, nil
}

func (ds *UserDatastore) GetByID(id string) (*v1.User, error) {
	user := &v1.User{}

	err := ds.db.Where("id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) GetByEmail(email string) (*v1.User, error) {
	user := &v1.User{}

	err := ds.db.Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) GetByVerificationCode(code string) (*v1.User, error) {
	user := &v1.User{}

	err := ds.db.Where("verification_code = ?", code).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by verification code: %s", err.Error())
	}

	return user, nil
}

func (ds *UserDatastore) Register(email, password string) (*v1.User, error) {
	tx := ds.db.Begin()

	user := &v1.User{}
	err := tx.Where("email = ?", email).First(user).Error
	if err == nil {
		tx.Rollback()
		return nil, ErrUserIsAlreadyExists
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

	passwordHash, _ := HashPassword(password)

	user = &v1.User{
		Id:        id,
		Email:     email,
		Password:  passwordHash,
		CreatedAt: pointer.ToTime(time.Now()),
	}

	err = tx.Create(user).Error
	if err != nil {
		ec := dbutil.ErrorCode(err)
		if ec == dbutil.ErrDuplicateEntry {
			return nil, ErrUserIsAlreadyExists
		}

		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return user, nil
}

func (ds *UserDatastore) UpdateAuthToken(user *v1.User, token string) error {
	user.Token = token

	updates := map[string]interface{}{
		"token": user.Token,
	}

	err := ds.db.Model(user).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}

func (ds *UserDatastore) ResetAuthToken(user *v1.User) error {
	user.Token = ""

	updates := map[string]interface{}{
		"token": user.Token,
	}

	err := ds.db.Model(user).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
