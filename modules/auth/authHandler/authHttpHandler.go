package authHandler

import (
	"context"
	"net/http"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/request"
	"github.com/chakornpat-tn/go-microservices/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	AuthHandlerService interface {
		Login(c echo.Context) error
	}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authUsecase.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authUsecase.AuthUsecaseService) AuthHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
	}
}

func (h *authHttpHandler) Login(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(auth.PlayerLoginReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.authUsecase.Login(ctx, h.cfg, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusOK, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
