package main

import (
	"lWords/config"
	"lWords/logger"
)

func main() {
	log := logger.Init()
	cfg := config.New()

	services := cfg.NewServices(log)

	//services.Migrations.Up(database.Username{
	//	ID:       primitive.NewObjectID(),
	//	Username: "test",
	//	Email:    "test",
	//	Password: "test",
	//})

	services.Migrations.CreateWords()
	//services.Migrations.CreateWordsIndexes()

	services.Migrations.Disconnect()
}
