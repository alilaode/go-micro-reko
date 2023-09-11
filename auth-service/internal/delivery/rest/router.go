package rest

import (
	"github.com/labstack/echo/v4"
)

func LoadRouters(e *echo.Echo, handler *handler) {

	userGroup := e.Group("/user")
	userGroup.POST("/register", handler.RegisterUser)
	userGroup.POST("/login", handler.Login)
}
