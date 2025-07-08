package server

import (
	"github.com/chakornpat-tn/go-microservices/modules/payment/paymentHandler"
	"github.com/chakornpat-tn/go-microservices/modules/payment/paymentRepository"
	"github.com/chakornpat-tn/go-microservices/modules/payment/paymentUsecase"
)

func (s *server) paymentService() {
	repo := paymentRepository.NewPaymentRepository(s.db)
	usecase := paymentUsecase.NewPaymentUsecase(repo)
	httpHandler := paymentHandler.NewPaymentHttpHandler(s.cfg, usecase)

	payment := s.app.Group("/payment_v1")

	// Health Check
	payment.GET("", s.healthCheckService)
	payment.POST("/payment/buy", httpHandler.BuyItem, s.middleware.JwtAuthorization)
	payment.POST("/payment/sell", httpHandler.SellItem, s.middleware.JwtAuthorization)

}
