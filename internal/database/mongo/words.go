package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"lWords/internal/models"
)

const (
	collectionName = "words"
)

func (m DBService) GetWordsList() ([]models.Words, error) {
	coll := m.client.Database(m.cfg.DBName).Collection(collectionName)

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		m.log.Error("failed to what??")
		return nil, err
	}

	var results []models.Words
	if err = cursor.All(context.TODO(), &results); err != nil {
		m.log.Error("failed to cursor words", zap.Error(err))
		return nil, err
	}

	return results, nil
}

func (m DBService) FindWords(word string) ([]models.Words, error) {
	coll := m.client.Database(m.cfg.DBName).Collection(collectionName)

	filter := bson.D{{"$text", bson.D{{"$search", word}}}}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		m.log.Error("failed to find words", zap.Error(err))
		return nil, err
	}

	var result []models.Words
	if err = cursor.All(context.Background(), &result); err != nil {
		m.log.Error("failed to cursor words", zap.Error(err))
		return nil, err
	}

	return result, nil
}
