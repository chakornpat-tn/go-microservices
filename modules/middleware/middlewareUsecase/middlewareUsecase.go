package middlewareUsecase

import "github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareRepository"

type (
	MiddlewareUsecaseService interface {
	}

	middlewareUsecase struct {
		middlewareRepo middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepo middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{
		middlewareRepo: middlewareRepo,
	}

}
