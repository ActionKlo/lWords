package migrations

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io"
	"lWords/internal/models"
	"net/http"
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

func (m *Migrations) Disconnect() {
	if err := m.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func (m *Migrations) CreateUser(user models.User) {
	coll := m.database.Collection("users")

	res, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		m.logger.Error("failed to insert user", zap.Error(err))
	}

	m.logger.Info(fmt.Sprintf("%s", res.InsertedID))
}

func (m *Migrations) FindUserByID(id primitive.ObjectID) {
	var user models.User
	err := m.database.Collection("").FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		m.logger.Error("failed to find user by id", zap.Error(err))
	}
	m.logger.Info(user.Username)
}

func (m *Migrations) CreateWords() {
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
		(*words)[i].Examples = []models.Sentences{}
		documents[i] = (*words)[i]
	}

	return documents, nil
}

func (m *Migrations) CreateWordsIndexes() {
	coll := m.database.Collection("words")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"eng", "text"},
			{"rus", "text"},
			{"ukr", "text"},
			{"pln", "text"},
		},
	}

	name, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}

	fmt.Println("name of index:", name)
}

func (m *Migrations) GenerateSentences() {
	fmt.Println("gtp request started")
	coll := m.database.Collection("words")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		m.logger.Fatal("failed to find words in db", zap.Error(err))
	}

	var words []models.Words
	if err = cursor.All(context.TODO(), &words); err != nil {
		m.logger.Fatal("failed to decode words", zap.Error(err))
	}

	for i := 0; i < len(words); i++ {
		sentences, err := postGPT(words[i], &i)
		if err != nil {
			m.logger.Fatal("filed to get data from GPT API", zap.Error(err))
		} else if sentences != nil {
			fmt.Println(sentences)
			err := coll.FindOneAndUpdate(
				context.TODO(),
				bson.D{{"_id", words[i].ID}},
				bson.D{{"$set", bson.D{{"examples", sentences}}}},
			)
			if err.Err() != nil {
				fmt.Println(err.Err())
				m.logger.Error(fmt.Sprintf("failed to update document, word id: %sentences", words[i].ID), zap.Error(err.Err()))
			}

			m.logger.Info("")
		}
	}

	return
}

func postGPT(word models.Words, i *int) ([]models.Sentences, error) {
	// TODO move to .env
	url := "http://100.104.232.63:1337/v1/chat/completions"

	body := models.GptBody{
		Model:  "gpt-3.5-turbo-16k",
		Stream: false,
	}
	body.Messages[0].Role = "assistant"
	body.Messages[0].Content = fmt.Sprintf("сгенерируй  3 предложения в виде JSON "+
		"[{ ru: text, en: text },{ ru: text, en: text },{ ru: text, en: text }] на русском и с переводом на "+
		"английский со словом %s, которое в контексте переводится на английский как %s. "+
		"Без лишних слов, без приветсвия, без форматировани и без лишних символов которые могут помешать декодированию сообщение на го "+
		"(у меня: json.NewDecoder(bytes.NewBuffer([]byte(r.Choices[0].Message.Content))).Decode(&s)); "+
		"только предложения в формате json массива", word.Rus, word.Eng)

	jsonBody, err := json.Marshal(body)

	if err != nil {
		fmt.Println("failed to marshal body:", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("failed to send request:", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("failed to close res body")
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("responce status code is: %d", res.StatusCode))
	}

	var r models.PostResponse
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("failed all:", err)
	}

	var s []models.Sentences

	if err = json.NewDecoder(bytes.NewBuffer([]byte(r.Choices[0].Message.Content))).Decode(&s); err != nil {
		fmt.Println("failed decode content:", err)
		*i--
		return nil, nil
	}

	return s, nil
}
