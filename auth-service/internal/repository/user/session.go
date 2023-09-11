package user

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/model"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

func (ur *UserRepo) CreateUserSession(ctx context.Context, userID string) (model.UserSession, error) {

	accessToken, err := ur.generateAccessToken(ctx, userID)
	if err != nil {
		return model.UserSession{}, err
	}

	return model.UserSession{
		JWTToken: accessToken,
	}, nil

}

func (ur *UserRepo) CheckSession(ctx context.Context, data model.UserSession) (userID string, err error) {

	accessToken, err := jwt.ParseWithClaims(data.JWTToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return &ur.signKey.PublicKey, nil
	})
	if err != nil {
		return "", err
	}

	// apakah claims bisa di masukkan ke tipe data *Claims
	accessTokenClaims, ok := accessToken.Claims.(*Claims)
	if !ok {
		return "", errors.New("unauthorized")
	}

	if accessToken.Valid {
		return accessTokenClaims.Subject, nil
	}

	return "", errors.New("unauthorized")
}

func (ur *UserRepo) generateAccessToken(ctx context.Context, userID string) (string, error) {

	accessTokenExp := time.Now().Add(2 * ur.accessExp).Unix()
	accessClaims := Claims{
		jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: accessTokenExp,
		},
	}

	accessJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), accessClaims)

	return accessJwt.SignedString(ur.signKey)

}
