package inventory

import (
	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/modules/models"
)

type (
	UpdateInventoryReq struct {
		PlayerID string `json:"player_id" validate:"required, max=64"`
		ItemID   string `json:"item_id" validate:"required, max=64"`
	}

	ItemInInventory struct {
		InventoryID string `json:"inventory_id"`
		PlayerID    string `json:"player_id"`
		*item.ItemShowCase
	}

	InventorySearchReq struct {
		models.PaginateReq
	}
)
