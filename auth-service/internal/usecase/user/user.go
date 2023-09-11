package user

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"context"
	"errors"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo user.Repository
}

func NewUseCase(userRepo user.Repository) *userUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (r *userUseCase) RegisterUser(ctx context.Context, request model.RegisterRequest) (model.User, error) {

	userRegistered, err := r.userRepo.CheckRegister(ctx, request.Username)
	if err != nil {
		return model.User{}, err
	}

	if userRegistered {
		return model.User{}, errors.New("user already register")
	}

	userHash, err := r.userRepo.GenerateUserHash(ctx, request.Password)
	if err != nil {
		return model.User{}, err
	}

	userData, err := r.userRepo.RegisterUser(ctx, model.User{
		ID:       uuid.NewString(),
		Username: request.Username,
		Hash:     userHash,
	})
	if err != nil {
		return model.User{}, err
	}

	return userData, nil

}

func (r *userUseCase) Login(ctx context.Context, request model.LoginRequest) (model.UserSession, error) {

	userData, err := r.userRepo.GetUserData(ctx, request.Username)
	if err != nil {
		return model.UserSession{}, err
	}

	verified, err := r.userRepo.VerifyLogin(ctx, request.Username, request.Password, userData)
	if err != nil {
		return model.UserSession{}, err
	}

	if !verified {
		return model.UserSession{}, errors.New("can't verify user login")
	}

	userSession, err := r.userRepo.CreateUserSession(ctx, userData.ID)
	if err != nil {
		return model.UserSession{}, err
	}

	return userSession, nil
}

func (r *userUseCase) CheckSession(ctx context.Context, sessionData model.UserSession) (userID string, err error) {

	userID, err = r.userRepo.CheckSession(ctx, sessionData)
	if err != nil {
		return "", err
	}

	return userID, err
}
