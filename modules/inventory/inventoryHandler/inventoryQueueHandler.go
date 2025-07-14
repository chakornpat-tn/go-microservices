package inventoryHandler

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/queue"
)

type (
	InventoryQueueHandlerService interface {
		AddPlayerItem()
		RemovePlayerItem()
		RollbackAddPlayerItem()
		RollbackRemovePlayerItem()
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase,
	}
}
func (h *inventoryQueueHandler) InventoryConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0", err.Error())
		consumer, err = worker.ConsumePartition("inventory", 0, 0)
		if err != nil {
			log.Printf("Error: InventoryConsumer failed:%s", err.Error())
			return nil, err
		}
	}
	return consumer, nil
}

func (h *inventoryQueueHandler) AddPlayerItem() {
	ctx := context.Background()
	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start AddPlayerItem....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerItem", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.AddPlayerItemRes(ctx, h.cfg, req)

				log.Printf("\n AddPlayerItem |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop AddPlayerItem")
			return
		}

	}
}
func (h *inventoryQueueHandler) RollbackAddPlayerItem() {
	ctx := context.Background()
	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start RollbackAddPlayerItem....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackAddPlayerItem", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "radd" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventory)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackAddPlayerItem(ctx, h.cfg, req)

				log.Printf("\n RollbackAddPlayerItem |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop RollbackAddPlayerItem")
			return
		}

	}
}
func (h *inventoryQueueHandler) RemovePlayerItem() {
	ctx := context.Background()
	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start RemovePlayerItem....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: RemovePlayerItem", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RemovePlayerItemRes(ctx, h.cfg, req)

				log.Printf("\n RemovePlayerItem |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop RemovePlayerItem")
			return
		}

	}
}
func (h *inventoryQueueHandler) RollbackRemovePlayerItem() {
	ctx := context.Background()
	consumer, err := h.InventoryConsumer(ctx)
	if err != nil {
		return
	}

	log.Println("Start RollbackRemovePlayerItem....")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	defer consumer.Close()

	for {

		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackRemovePlayerItem", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rremove" {
				h.inventoryUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventory)
				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackRemovePlayerItem(ctx, h.cfg, req)

				log.Printf("\n RollbackRemovePlayerItem |Topic[%s] | Partition[%d]] | Offset[%d] | Message[%s] \n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
		case <-sigChan:
			log.Println("Stop RollbackRemovePlayerItem")
			return
		}

	}
}
