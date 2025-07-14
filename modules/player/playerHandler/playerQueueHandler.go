package playerHandler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/queue"
)

type (
	PlayerQueueHandlerService interface {
		DockedPlayerMoney()
		RollBackPlayerTransaction()
		AddPlayerMoney()
	}

	playerQueueHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerQueueHanddler(cfg *config.Config, playerUsecase playerUsecase.PlayerUsecaseService) PlayerQueueHandlerService {
	return &playerQueueHandler{
		cfg:           cfg,
		playerUsecase: playerUsecase,
	}
}

func (h *playerQueueHandler) PlayerConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.playerUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("player", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0", err.Error())
		consumer, err = worker.ConsumePartition("player", 0, 0)
		if err != nil {
			return nil, err
		}
	}
	return consumer, nil
}

func (h *playerQueueHandler) DockedPlayerMoney() {
	ctx := context.Background()
	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start DockedPlayerMoney....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: DockedPlayerMoney", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.DockedPlayerMoneyRes(ctx, h.cfg, req)

				log.Printf("\n DockedPlayerMoney |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop DockedPlayerMoney")
			return
		}

	}
}

func (h *playerQueueHandler) RollBackPlayerTransaction() {
	ctx := context.Background()
	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start RollBackPlayerTransaction....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollBackPlayerTransaction", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rtransaction" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.RollBackPlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.RollBackPlayerTransaction(ctx, req)

				log.Printf("\n RollBackPlayerTransaction |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop DockedPlayerMoney")
			return
		}

	}
}

func (h *playerQueueHandler) AddPlayerMoney() {
	ctx := context.Background()
	consumer, err := h.PlayerConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start AddPlayerMoney....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerMoney", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.playerUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.AddPlayerMoneyRes(ctx, h.cfg, req)

				log.Printf("\n AddPlayerMoney |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop AddPlayerMoney")
			return
		}

	}
}
