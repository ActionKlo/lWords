package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Config struct {
	Mongo
}

type Mongo struct {
	Uri      string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type MongoDBService struct {
	client *mongo.Client
	cfg    *Config
	log    *zap.Logger
}

func NewMongoDB(log *zap.Logger, cfg *Config) *MongoDBService {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.Mongo.Username,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("failed to connect to MongoDBService", zap.Error(err))
	}

	return &MongoDBService{
		client: client,
		cfg:    cfg,
		log:    log,
	}
}
