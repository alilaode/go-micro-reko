package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"auth-service/internal/model"
	"auth-service/internal/model/constant"
	"auth-service/internal/usecase/user"
	"auth-service/internal/utils/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type authMiddleware struct {
	userUseCase user.UseCase
}

func Init(e *echo.Echo) {
	e.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}))
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
}

func GetAuthMiddleware(userUseCase user.UseCase) *authMiddleware {
	return &authMiddleware{
		userUseCase: userUseCase,
	}
}

func (am *authMiddleware) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sessionData, err := GetSessionData(c.Request())
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}

		userID, err := am.userUseCase.CheckSession(c.Request().Context(), sessionData)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  err.Error(),
				Internal: err,
			}
		}

		//buat context baru
		authContext := context.WithValue(c.Request().Context(), constant.AuthContextKey, userID)
		c.SetRequest(c.Request().WithContext(authContext))

		if err := next(c); err != nil {
			return err
		}

		return nil
	}
}

func GetSessionData(r *http.Request) (model.UserSession, error) {
	authString := r.Header.Get("Authorization")
	splitString := strings.Split(authString, " ")
	if len(splitString) != 2 {
		return model.UserSession{}, errors.New("unauthorization")
	}
	accessString := splitString[1]

	return model.UserSession{
		JWTToken: accessString,
	}, nil
}
