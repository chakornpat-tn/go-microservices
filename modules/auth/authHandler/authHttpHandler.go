package authHandler

import (
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
)

type (
	AuthHandlerService interface {
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authUsecase.AuthUsecaseService) AuthHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
	}
}
