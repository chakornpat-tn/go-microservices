package playerHandler

import "github.com/chakornpat-tn/go-microservices/modules/player/playerUsecase"

type (
	playerGrpcHandler struct {
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) playerUsecase.PlayerUsecaseService {
	return &playerGrpcHandler{
		playerUsecase: playerUsecase,
	}
}
