package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"time"

	"auth-service/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db        *gorm.DB
	gcm       cipher.AEAD
	time      uint32
	memory    uint32
	threads   uint8
	keylen    uint32
	signKey   *rsa.PrivateKey
	accessExp time.Duration
}

func NewRepository(
	db *gorm.DB,
	secret string,
	time uint32,
	memory uint32,
	threads uint8,
	keylen uint32,
	signKey *rsa.PrivateKey,
	accessExp time.Duration,
) (Repository, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		db:        db,
		gcm:       gcm,
		time:      time,
		memory:    memory,
		threads:   threads,
		keylen:    keylen,
		signKey:   signKey,
		accessExp: accessExp,
	}, nil
}

func (ur *UserRepo) RegisterUser(ctx context.Context, userData model.User) (model.User, error) {

	if err := ur.db.WithContext(ctx).Create(&userData).Error; err != nil {
		return model.User{}, err
	}
	return userData, nil
}

func (ur *UserRepo) CheckRegister(ctx context.Context, username string) (bool, error) {
	var userData model.User

	if err := ur.db.WithContext(ctx).Where(model.User{Username: username}).Take(&userData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	return userData.ID != "", nil
}

func (ur *UserRepo) GetUserData(ctx context.Context, username string) (userData model.User, err error) {

	err = ur.db.WithContext(ctx).Where(model.User{Username: username}).Take(&userData).Error

	return userData, err
}

func (ur *UserRepo) VerifyLogin(ctx context.Context, username, password string, userData model.User) (bool, error) {

	if username != userData.Username {
		return false, nil
	}

	verified, err := ur.comparePassword(ctx, password, userData.Hash)
	if err != nil {
		return false, nil
	}

	return verified, nil

}
