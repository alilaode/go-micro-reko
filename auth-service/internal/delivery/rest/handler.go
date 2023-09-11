package rest

import "auth-service/internal/usecase/user"

type handler struct {
	userUseCase user.UseCase
}

func NewHandler(authUseCase user.UseCase) *handler {
	return &handler{userUseCase: authUseCase}
}
