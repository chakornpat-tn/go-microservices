package middlewareHandler

import (
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareUsecase"
)

type (
	MiddlewareHandlerService interface {
	}

	middlewareHandler struct {
		middlewareUsecase middlewareUsecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareUsecase middlewareUsecase.MiddlewareUsecaseService) MiddlewareHandlerService {
	return &middlewareHandler{
		middlewareUsecase: middlewareUsecase,
	}

}
