package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"auth-service/internal/model"
	"auth-service/internal/utils/response"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterUser(c echo.Context) error {

	var request model.RegisterRequest

	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	userData, err := h.userUseCase.RegisterUser(c.Request().Context(), request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": userData,
	})
}

func (h *Handler) Login(c echo.Context) error {

	var (
		err     error
		request = new(model.LoginRequest)
	)

	if err = c.Bind(request); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err = c.Validate(request); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	sessionData, err := h.userUseCase.Login(c.Request().Context(), *request)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("[delivery][rest][user_handler][Login] unable to login")
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(sessionData).Send(c)

}
