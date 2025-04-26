package server

import (
	"github.com/chakornpat-tn/go-microservices/modules/player/playerHandler"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerRepository"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerUsecase"
)

func (s *server) playerService() {
	repo := playerRepository.NewPlayerRepository(s.db)
	usecase := playerUsecase.NewPlayerUsecase(repo)
	httpHandler := playerHandler.NewPlayerHttpHanddler(s.cfg, usecase)
	grpcHandler := playerHandler.NewPlayerGrpcHandler(usecase)
	queueHandler := playerHandler.NewPlayerQueueHanddler(s.cfg, usecase)

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	player := s.app.Group("/player_v1")

	// Health Check
	player.GET("", s.healthCheckService)

}
