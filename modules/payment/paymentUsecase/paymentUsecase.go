package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/chakornpat-tn/go-microservices/modules/item"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/modules/payment"
	"github.com/chakornpat-tn/go-microservices/modules/payment/paymentRepository"
)

type (
	PaymentUsecaseService interface {
		FindItemsInIDs(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
	}

	paymentUsecase struct {
		paymentRepo paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepo paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
	}
}

func (u *paymentUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.paymentRepo.GetOffset(pctx)
}

func (u *paymentUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.paymentRepo.UpserOffset(pctx, offset)
}

func (u *paymentUsecase) FindItemsInIDs(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error {
	setIDs := make(map[string]bool)
	for _, v := range req {
		if !setIDs[v.ItemID] {
			setIDs[v.ItemID] = true
		}
	}

	itemData, err := u.paymentRepo.FindItemsInIDs(pctx, grpcUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for k := range setIDs {
				itemIds = append(itemIds, k)
			}
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: FindItemsInIDs failed: %s", err.Error())
		return errors.New("error: item not found")
	}

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemID:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImgUrl,
			Damage:   int(v.Damage),
		}
	}

	for i := range req {
		if _, ok := itemMaps[req[i].ItemID]; !ok {
			log.Printf("Error: item not found: %s", req[i].ItemID)
			return errors.New("error: item not found")
		}
		req[i].Price = itemMaps[req[i].ItemID].Price
	}

	return nil
}
