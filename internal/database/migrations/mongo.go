package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io"
	"lWords/internal/models"
	"os"
	"time"
)

type (
	Config struct {
		Mongo Mongo
	}

	Mongo struct {
		URI      string
		Host     string
		Port     string
		DBName   string
		User     string
		Password string
	}

	Migrations struct {
		database *mongo.Database
		client   *mongo.Client
		logger   *zap.Logger
		config   *Config
	}
)

func Init(logger *zap.Logger, cfg *Config) *Migrations {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		cfg.Mongo.DBName)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &Migrations{
		database: client.Database(cfg.Mongo.DBName),
		client:   client,
		logger:   logger,
		config:   cfg,
	}
}

func (m Migrations) Disconnect() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func (m Migrations) CreateUser(user models.User) {
	coll := m.database.Collection("users")

	res, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		m.logger.Error("failed to insert user", zap.Error(err))
	}

	m.logger.Info(fmt.Sprintf("%s", res.InsertedID))
}

func (m Migrations) FindUserByID(id primitive.ObjectID) {
	var user models.User
	err := m.database.Collection("").FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		m.logger.Error("failed to find user by id", zap.Error(err))
	}
	m.logger.Info(user.Username)
}

func (m Migrations) CreateWords() {
	err := m.database.Collection("words").Drop(context.Background())
	if err != nil {
		m.logger.Fatal("failed to drop words", zap.Error(err))
	}

	err = m.database.Collection("words").Drop(context.Background())
	if err != nil {
		m.logger.Fatal("failed to drop words", zap.Error(err))
	}

	start := time.Now()

	words, err := getWordsFromFile()
	if err != nil {
		m.logger.Fatal("failed to get words from file", zap.Error(err))
	}

	documents, err := prepareWordsObjects(&words)
	if err != nil {
		m.logger.Fatal("failed to make words objects", zap.Error(err))
	}

	coll := m.database.Collection("words")

	_, err = coll.InsertMany(context.Background(), documents)
	if err != nil {
		m.logger.Fatal("failed to insert words", zap.Error(err))
	}

	fmt.Println(time.Since(start))
}

func getWordsFromFile() ([]models.Words, error) {
	jsonFilePath := "internal/database/migrations/words.json"

	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {
			return
		}
	}(jsonFile)

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var words []models.Words
	if err = json.Unmarshal(data, &words); err != nil {
		return nil, err
	}

	return words, err
}

func prepareWordsObjects(words *[]models.Words) ([]interface{}, error) {
	documents := make([]interface{}, len(*words))

	for i := range *words {
		(*words)[i].ID = primitive.NewObjectID()
		documents[i] = (*words)[i]
	}

	return documents, nil
}

func (m Migrations) CreateWordsIndexes() {
	coll := m.database.Collection("words")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"eng", "text"},
			{"rus", "text"},
			{"urk", "text"},
			{"pln", "text"},
		},
	}

	name, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}

	fmt.Println("name of index:", name)
}
