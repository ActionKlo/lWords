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
	// TODO check if base/collections exist or make migrations from .sh file
	// migration.IfNotExist
	migration.CreateWords()
	migration.CreateWordsIndexes()
	migration.Disconnect()

	mongodb := services.MongoDB // TODO make services "web" "bot" "itd"
	_ = mongodb
}
