package handlers

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"lWords/internal/database/mongo"
	"net/http"
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

func (c *Config) GetAllWords(ctx echo.Context) error {
	words, err := c.mongo.GetWordsList()
	if err != nil {
		c.logger.Error("failed to get words list", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, words)
}

func (c *Config) FindWords(ctx echo.Context) error {
	words, err := c.mongo.FindWords(ctx.Param("word"))
	if err != nil {
		c.logger.Error("failed to get word", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, words)
}
