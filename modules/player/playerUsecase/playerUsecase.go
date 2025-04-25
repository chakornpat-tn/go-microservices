package playerUsecase

import "github.com/chakornpat-tn/go-microservices/modules/player/playerRepository"

type (
	PlayerUsecaseService interface{}

	playerUsecase struct {
		playerRepo playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepo playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{
		playerRepo: playerRepo,
	}
}
