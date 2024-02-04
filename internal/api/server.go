package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"lWords/internal/api/handlers"
	"lWords/internal/database/mongo"
)

type Config struct {
	Host string
	Port string
}

type EchoServer struct {
	logger *zap.Logger
	cfg    *Config
	mongo  *mongo.DBService
}

func NewEchoServer(logger *zap.Logger, mongoDBService *mongo.DBService, cfg *Config) *EchoServer {
	return &EchoServer{
		logger: logger,
		cfg:    cfg,
		mongo:  mongoDBService,
	}
}

func (s *EchoServer) Start() {
	app := echo.New()

	h := handlers.NewHandlerConfig(s.mongo, s.logger)

	app.GET("/words", h.GetAllWords)
	app.GET("/words/:word", h.FindWords)

	if err := app.Start(fmt.Sprintf(":%s", s.cfg.Port)); err != nil {
		s.logger.Fatal(fmt.Sprintf("failed to start start server on address %s:%s", s.cfg.Host, s.cfg.Port), zap.Error(err))
	}
}
