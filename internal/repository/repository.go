package repository

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	mongoDB       *mongo.Database
	redisClient   *redis.Client
	Authorization *AuthorizationRepository
	Admin         *AdminRepository
	Cash          *CashRepository
}

func NewRepository(mongoDB *mongo.Database, redisClient *redis.Client) *Repository {
	return &Repository{
		mongoDB:       mongoDB,
		redisClient:   redisClient,
		Authorization: newAuthorizationRepositoryMongoDB(mongoDB, redisClient),
		Admin:         newAdminRepository(mongoDB, redisClient),
		Cash:          newCashRepository(redisClient),
	}
}
