package inventoryHandler

import (
	"context"
	"net/http"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/request"
	"github.com/chakornpat-tn/go-microservices/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	InventoryHttpHandlerService interface {
		FindPlayerItems(c echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryHttpHandlerService {
	return &inventoryHttpHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryHttpHandler) FindPlayerItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerID := c.Param("player_id")

	req := new(inventory.InventorySearchReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.inventoryUsecase.FindPlayerItems(ctx, h.cfg.Paginate.InventoryNextPageBasedUrl, playerID, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, result)

}
