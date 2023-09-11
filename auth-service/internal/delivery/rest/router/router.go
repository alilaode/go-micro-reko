package router

import (
	"auth-service/internal/delivery/rest"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, handler *rest.Handler) {

	userGroup := e.Group("/user")
	userGroup.POST("/register", handler.RegisterUser)
	userGroup.POST("/login", handler.Login)
}
