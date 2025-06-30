package itemHandler

import (
	"context"
	"net/http"
	"strings"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/modules/item/itemUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/request"
	"github.com/chakornpat-tn/go-microservices/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	ItemHttpHandlerService interface {
		CreateItem(c echo.Context) error
		FindOneItem(c echo.Context) error
		FindManyItem(c echo.Context) error
		UpdateItem(c echo.Context) error
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecase.ItemUsecaseService) ItemHttpHandlerService {
	return &itemHttpHandler{
		cfg:         cfg,
		itemUsecase: itemUsecase,
	}
}

func (h *itemHttpHandler) CreateItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.CreateItem(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)

}
func (h *itemHttpHandler) FindOneItem(c echo.Context) error {
	ctx := context.Background()

	itemID := strings.TrimPrefix(c.Param("item_id"), "item:")

	result, err := h.itemUsecase.FindOneItem(ctx, itemID)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, result)

}

func (h *itemHttpHandler) FindManyItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemSearchReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.itemUsecase.FindManyItems(ctx, h.cfg.Paginate.ItemNextPageBasedUrl, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, result)
}

func (h *itemHttpHandler) UpdateItem(c echo.Context) error {
	ctx := context.Background()

	itemID := strings.TrimPrefix(c.Param("item_id"), "item:")

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemUpdateReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.EditItem(ctx, itemID, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
