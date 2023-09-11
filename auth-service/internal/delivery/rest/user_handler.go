package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"auth-service/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *handler) RegisterUser(c echo.Context) error {
	/*
		ctx, span := tracing.CreateSpan(c.Request().Context(), "RegisterUser")
		defer span.End()
	*/
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

func (h *handler) Login(c echo.Context) error {
	var request model.LoginRequest

	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	sessionData, err := h.userUseCase.Login(c.Request().Context(), request)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("[delivery][rest][user_handler][Login] unable to login")

		//fmt.Printf("got error %s\n", err.Error())

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": sessionData,
	})

}
