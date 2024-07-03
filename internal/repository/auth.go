package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"petProject/internal/model"
	"time"
)

type AuthorizationRepository struct {
	collectionMongo *mongo.Collection
	redisClient     *redis.Client
}

func newAuthorizationRepositoryMongoDB(mongoDB *mongo.Database, redisClient *redis.Client) *AuthorizationRepository {
	return &AuthorizationRepository{
		collectionMongo: mongoDB.Collection("users"),
		redisClient:     redisClient,
	}
}

func (r *AuthorizationRepository) CreateUser(user model.User) (model.User, error) {
	ctx := context.TODO()

	// Input in mongoDB
	result, err := r.collectionMongo.InsertOne(ctx, user)

	if err != nil {
		logrus.Error("Failed to create user:", err)
		return user, err
	}

	// Отримання ідентифікатора нового користувача
	insertedID := result.InsertedID.(primitive.ObjectID)
	user.ID = insertedID
	logrus.Info("User created successfully")
	return user, nil
}

func (r *AuthorizationRepository) GetUser(name, password string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Створення фільтру для пошуку користувача
	filter := bson.M{
		"name":     name,
		"password": password,
	}

	// Виконання запиту FindOne з фільтром
	var user model.User
	err := r.collectionMongo.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Якщо користувача не знайдено
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthorizationRepository) CheckIsUserExist(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"name": username,
	}

	var user model.User
	err := r.collectionMongo.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return fmt.Errorf("User named %s already exists", username)
}

func (r *AuthorizationRepository) WriteTokenInRedis(username, token, device string) error {
	ctx := context.Background()
	key := username + ":" + device
	err := r.redisClient.Set(ctx, key, token, 15*time.Minute).Err() // 0 означає, що токен не буде мати терміну придатності
	if err != nil {
		return fmt.Errorf("failed to write token to Redis: %w", err)
	}
	return nil
}
