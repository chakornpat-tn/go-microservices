package itemRepository

import "go.mongodb.org/mongo-driver/v2/mongo"

type (
	ItemRepositoryService interface{}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{
		db: db,
	}
}
