package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"lWords/internal/api"
	"lWords/internal/database/migrations"
	"lWords/internal/database/mongo"
	"log"
)

type (
	Config struct {
		Mongo
		WebAPI
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

	WebAPI struct {
		Host string `mapstructure:"WEB_HOST"`
		Port string `mapstructure:"WEB_PORT"`
	}
)

type Services struct {
	Migrations *migrations.Migrations
	MongoDB    *mongo.DBService
	WebAPI     *api.EchoServer
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

	if err := v.Unmarshal(&appConfig.WebAPI); err != nil {
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
		Uri:      c.Mongo.URI,
		Host:     c.Mongo.Host,
		Port:     c.Mongo.Port,
		Username: c.Mongo.Username,
		Password: c.Mongo.Password,
		DBName:   c.Mongo.DBName,
	}})

	webAPI := api.NewEchoServer(logger, mongodb, &api.Config{
		Host: c.WebAPI.Host,
		Port: c.WebAPI.Port,
	})

	return &Services{
		Migrations: mongomig,
		MongoDB:    mongodb,
		WebAPI:     webAPI,
	}
}
