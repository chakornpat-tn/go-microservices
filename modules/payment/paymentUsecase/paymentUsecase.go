package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	"github.com/chakornpat-tn/go-microservices/modules/item"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/modules/payment"
	"github.com/chakornpat-tn/go-microservices/modules/payment/paymentRepository"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/pkg/queue"
)

type (
	PaymentUsecaseService interface {
		FindItemsInIDs(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		BuyItem(pctx context.Context, cfg *config.Config, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		SellItem(pctx context.Context, cfg *config.Config, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error)
		BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes)
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

func (u *paymentUsecase) PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{cfg.Kafka.Url}, cfg.Kafka.ApiKey, cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepo.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0", err.Error())
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			return nil, err
		}
	}
	return consumer, nil
}

func (u *paymentUsecase) BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx, cfg)
	if err != nil {
		resCh <- nil
		return
	}

	defer consumer.Close()

	log.Println("Start BuyOrSellConsumer....")

	select {
	case err := <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpserOffset(pctx, msg.Offset+1)
			req := new(payment.PaymentTransferRes)
			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("\n BuyOrSellConsumer |Topic[%s] | Partition[%d] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
		}
	}

}

func (u *paymentUsecase) BuyItem(pctx context.Context, cfg *config.Config, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIDs(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepo.DockedPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerID: playerID,
			Amount:   -item.Price,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)
		res := <-resCh
		if res != nil {
			log.Println(res)
			stage1 = append(stage1, &payment.PaymentTransferRes{
				TransactionID: res.TransactionID,
				InventoryID:   "",
				PlayerID:      playerID,
				ItemID:        item.ItemID,
				Amount:        item.Price,
				Error:         res.Error,
			})
		}

		for _, s1 := range stage1 {
			if s1.Error != "" {
				for _, ss1 := range stage1 {
					u.paymentRepo.RollbackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
						TransactionID: ss1.TransactionID,
					})
				}
				return nil, errors.New("error: buy item failed")
			}
		}

	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		u.paymentRepo.AddPlayerItem(pctx, cfg, &inventory.UpdateInventoryReq{
			PlayerID: s1.PlayerID,
			ItemID:   s1.ItemID,
		})

		resCh := make(chan *payment.PaymentTransferRes)
		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage2 = append(stage2, &payment.PaymentTransferRes{
				TransactionID: s1.TransactionID,
				InventoryID:   res.InventoryID,
				PlayerID:      playerID,
				ItemID:        s1.ItemID,
				Amount:        s1.Amount,
				Error:         s1.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {
			for _, ss2 := range stage2 {
				u.paymentRepo.RollbackAddPlayerItem(pctx, cfg, &inventory.RollbackPlayerInventory{
					InventoryID: ss2.InventoryID,
				})
			}

			for _, ss2 := range stage2 {
				u.paymentRepo.RollbackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
					TransactionID: ss2.TransactionID,
				})
			}

			return nil, errors.New("error: buy item failed")
		}

	}

	return stage2, nil
}

func (u *paymentUsecase) SellItem(pctx context.Context, cfg *config.Config, playerID string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIDs(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepo.RemovePlayerItem(pctx, cfg, &inventory.UpdateInventoryReq{
			PlayerID: playerID,
			ItemID:   item.ItemID,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "sell", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage1 = append(stage1, &payment.PaymentTransferRes{
				InventoryID:   "",
				TransactionID: "",
				PlayerID:      playerID,
				ItemID:        item.ItemID,
				Amount:        item.Price,
				Error:         res.Error,
			})
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				if ss1.Error != "error: item not found" {
					u.paymentRepo.RollbackRemovePlayerItem(pctx, cfg, &inventory.RollbackPlayerInventory{
						PlayerID: playerID,
						ItemID:   ss1.ItemID,
					})
				}
			}
			return nil, errors.New("error: sell item failed")
		}
	}

	stage2 := make([]*payment.PaymentTransferRes, 0)
	for _, s1 := range stage1 {
		u.paymentRepo.AddPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerID: playerID,
			Amount:   s1.Amount * 0.5,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "sell", cfg, resCh)

		res := <-resCh
		if res != nil {
			log.Println(res)
			stage2 = append(stage2, &payment.PaymentTransferRes{
				InventoryID:   "",
				TransactionID: s1.TransactionID,
				PlayerID:      playerID,
				ItemID:        s1.ItemID,
				Amount:        s1.Amount,
				Error:         s1.Error,
			})
		}
	}

	for _, s2 := range stage2 {
		if s2.Error != "" {
			for _, ss2 := range stage2 {
				u.paymentRepo.RollbackTransaction(pctx, cfg, &player.RollBackPlayerTransactionReq{
					TransactionID: ss2.TransactionID,
				})
			}

			for _, ss2 := range stage2 {
				if ss2.Error != "error: item not found" {
					u.paymentRepo.RollbackRemovePlayerItem(pctx, cfg, &inventory.RollbackPlayerInventory{
						InventoryID: ss2.InventoryID,
					})
				}
			}

			return nil, errors.New("error: sell item failed")
		}
	}

	return stage2, nil
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
