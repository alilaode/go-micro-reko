package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	logger "auth-service/internal/Logger"
	"auth-service/internal/database"
	"auth-service/internal/delivery/rest"

	uRepo "auth-service/internal/repository/user"
	uUseCase "auth-service/internal/usecase/user"

	"github.com/labstack/echo/v4"
)

const (
	//DSN = "host=postgres port=5454 user=postgres password=root dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	webPort = 80
)

func main() {

	e := echo.New()

	logger.Init()

	db := database.GetDB(os.Getenv("DSN"))

	secret := "AES256Key-32Characters1234567890"

	signKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	userRepo, err := uRepo.NewRepository(db, secret, 1, 64*1024, 4, 32, signKey, 60*time.Second)
	if err != nil {
		panic(err)
	}
	userUseCase := uUseCase.NewUseCase(userRepo)

	h := rest.NewHandler(userUseCase)

	rest.LoadMiddlewares(e)
	rest.LoadRouters(e, h)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", webPort)))

}
