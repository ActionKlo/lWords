package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"lWords/database/migrations"
	"lWords/database/mongo"
	"log"
)

type (
	Config struct {
		Mongo
	}

	Mongo struct {
		DBUrl    string `mapstructure:"DB_URL"`
		URI      string `mapstructure:"MONGO_URI"`
		Username string `mapstructure:"MONGO_USER"`
		Password string `mapstructure:"MONGO_PASSWORD"`
		Host     string `mapstructure:"MONGO_HOST"`
		Port     string `mapstructure:"MONGO_PORT"`
		DBName   string `mapstructure:"MONGO_DB_NAME"`
	}
)

type Services struct {
	Migrations *migrations.Migrations
	MongoDB    *mongo.MongoDBService
}

func New() *Config {
	var appConfig Config
	v := viper.New()
	v.SetConfigType("dotenv")
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := v.Unmarshal(&appConfig.Mongo); err != nil {
		log.Fatal(err)
	}
	return &appConfig
}

func (c *Config) NewServices(logger *zap.Logger) *Services {
	mongomig := migrations.Init(logger, &migrations.Config{
		Mongo: migrations.Mongo{
			User:     c.Mongo.Username,
			Password: c.Mongo.Password,
			URI:      c.Mongo.URI,
			Host:     c.Mongo.Host,
			Port:     c.Mongo.Port,
			DBName:   c.Mongo.DBName,
		},
	})

	mongodb := mongo.NewMongoDB(logger, &mongo.Config{Mongo: mongo.Mongo{
		Uri:      c.URI,
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		DBName:   c.DBName,
	}})

	return &Services{
		Migrations: mongomig,
		MongoDB:    mongodb,
	}
}
