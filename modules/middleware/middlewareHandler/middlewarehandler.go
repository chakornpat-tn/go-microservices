package middlewareHandler

import "github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareUsecase"

type (
	middlewareHandlerService interface {
	}

	middlewareHandler struct {
		middlewareUsecase middlewareUsecase.MiddlewareUsecaseService
	}
)

func NewMiddlewareHandler(middlewareUsecase middlewareUsecase.MiddlewareUsecaseService) middlewareHandlerService {
	return &middlewareHandler{
		middlewareUsecase: middlewareUsecase,
	}

}
