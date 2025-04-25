package playerRepository

import "go.mongodb.org/mongo-driver/v2/mongo"

type (
	PlayerRepositoryService interface {
	}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayerRepository(db *mongo.Client) PlayerRepositoryService {
	return &playerRepository{
		db: db,
	}
}
