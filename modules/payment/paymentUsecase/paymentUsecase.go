package paymentUsecase

import "github.com/chakornpat-tn/go-microservices/modules/payment/paymentRepository"

type (
	PaymentUsecaseService interface{}

	paymentUsecase struct {
		paymentRepo paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepo paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
	}
}
