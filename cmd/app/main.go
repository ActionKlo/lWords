package main

import (
	"lWords/config"
	"lWords/logger"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	services := cfg.NewServices(log)

	migration := services.Migrations
	migration.CreateWords()
	//TODO add words indexes
	migration.Disconnect()

	mongodb := services.MongoDB // TODO make services "web" "bot" "itd"
	_ = mongodb
}
