package user

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/user"
	"auth-service/internal/utils/response"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	/*
		userRegistered, err := r.userRepo.CheckRegister(ctx, request.Username)
		if err != nil {
			return model.User{}, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}

		if userRegistered {
			return model.User{}, response.ErrorBuilder(&response.ErrorConstant.Duplicate, errors.New("user already register"))
		}
	*/

	userHash, err := r.userRepo.GenerateUserHash(ctx, request.Password)
	if err != nil {
		return model.User{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	userData, err := r.userRepo.RegisterUser(ctx, model.User{
		ID:       uuid.NewString(),
		Username: request.Username,
		Hash:     userHash,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.User{}, response.ErrorBuilder(&response.ErrorConstant.Duplicate, err)
		}
		return model.User{}, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	return userData, nil

}

func (r *userUseCase) Login(ctx context.Context, request model.LoginRequest) (model.UserSession, error) {

	userData, err := r.userRepo.GetUserData(ctx, request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserSession{}, response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err)
		}
		return model.UserSession{}, err
	}

	verified, err := r.userRepo.VerifyLogin(ctx, request.Username, request.Password, userData)
	if err != nil {
		return model.UserSession{}, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	if !verified {
		return model.UserSession{}, response.ErrorBuilder(&response.ErrorConstant.Unauthorized, errors.New("can't verify user login"))
	}

	userSession, err := r.userRepo.CreateUserSession(ctx, userData.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.UserSession{}, response.ErrorBuilder(&response.ErrorConstant.Duplicate, err)
		}
		return model.UserSession{}, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	return userSession, nil
}

func (r *userUseCase) CheckSession(ctx context.Context, sessionData model.UserSession) (userID string, err error) {

	userID, err = r.userRepo.CheckSession(ctx, sessionData)
	if err != nil {
		return "", response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	return userID, err
}
