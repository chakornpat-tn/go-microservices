package playerHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/request"
	"github.com/chakornpat-tn/go-microservices/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	PlayerHttpHandlerService interface {
		CreatePlayer(c echo.Context) error
		FindOnePlayer(c echo.Context) error
		AddPlayerMoney(c echo.Context) error
	}

	playerHttpHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerHttpHanddler(cfg *config.Config, playerUsecase playerUsecase.PlayerUsecaseService) PlayerHttpHandlerService {
	return &playerHttpHandler{
		cfg:           cfg,
		playerUsecase: playerUsecase,
	}
}

func (h *playerHttpHandler) CreatePlayer(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(player.CreatePlayerReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.playerUsecase.CreatePlayer(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) FindOnePlayer(c echo.Context) error {
	ctx := context.Background()
	playerID := strings.TrimPrefix(c.Param("playerID"), "player:")

	res, err := h.playerUsecase.FindOnePlayer(ctx, playerID)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *playerHttpHandler) AddPlayerMoney(c echo.Context) error {
	ctx := context.Background()
	wrapper := request.ContextWrapper(c)

	req := new(player.CreatePlayerTransactionReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.playerUsecase.AddPlayerMonney(ctx, req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{"message": "successful"})
}
