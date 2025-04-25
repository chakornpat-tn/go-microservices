package authUsecase

import "github.com/chakornpat-tn/go-microservices/modules/auth/authRepository"

type (
	AuthUsecaseService interface {
	}

	authUsecase struct {
		authRepo authRepository.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepo authRepository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{
		authRepo: authRepo,
	}
}
