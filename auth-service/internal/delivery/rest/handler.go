package rest

import "auth-service/internal/usecase/user"

type Handler struct {
	userUseCase user.UseCase
}

func NewHandler(authUseCase user.UseCase) *Handler {
	return &Handler{userUseCase: authUseCase}
}
