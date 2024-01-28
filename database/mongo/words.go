package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"lWords/models"
)

func (m MongoDBService) FindWords(word string) []models.Words {
	coll := m.client.Database(m.cfg.DBName).Collection("words")

	filter := bson.D{{"$text", bson.D{{"$search", word}}}}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		m.log.Error("failed to find words", zap.Error(err))
	}

	var result []models.Words
	if err = cursor.All(context.Background(), &result); err != nil {
		m.log.Error("failed to cursor words", zap.Error(err))
	}

	return result
}

func (m MongoDBService) FindWordByID() []models.Words {
	coll := m.client.Database(m.cfg.DBName).Collection("words")

	objID, err := primitive.ObjectIDFromHex("65abb6fe0db988e11c8a6736")
	if err != nil {
		m.log.Error("failed to create ObjectID from hex", zap.Error(err))
	}
	fmt.Println(objID)

	filter := bson.D{{"_id", objID}}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		m.log.Error("failed to find word by ID", zap.Error(err))
	}

	var result []models.Words
	if err = cursor.All(context.Background(), &result); err != nil {
		m.log.Error("failed to cursor word", zap.Error(err))
	}

	return result
}
