package playerHandler

import "github.com/chakornpat-tn/go-microservices/modules/player/playerUsecase"

type (
	playerGrpcHandler struct {
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) *playerGrpcHandler {
	return &playerGrpcHandler{
		playerUsecase: playerUsecase,
	}
}
