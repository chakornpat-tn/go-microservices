package authHandler

import (
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
)

type (
	AuthHandlerService interface {
	}

	authHttpHandler struct {
		authUsecase authUsecase.AuthUsecaseService
	}
)

func NewAuthHandler(cfg *config.Config, authUsecase authUsecase.AuthUsecaseService) AuthHandlerService {
	return &authHttpHandler{
		authUsecase: authUsecase,
	}
}
