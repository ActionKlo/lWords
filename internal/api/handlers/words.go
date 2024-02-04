package handlers

import (
	"go.uber.org/zap"
	"lWords/internal/database/mongo"
)

type Config struct {
	mongo  *mongo.DBService
	logger *zap.Logger
}

func NewHandlerConfig(mongo *mongo.DBService, logger *zap.Logger) *Config {
	return &Config{
		mongo:  mongo,
		logger: logger,
	}
}
