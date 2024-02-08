package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	StudyResultID primitive.ObjectID `bson:"studyResultID"`
	StatisticsID  primitive.ObjectID `bson:"statisticsID"`
	Username      string             `bson:"username"`
	Email         string             `bson:"email"`
	Password      string             `bson:"password"`
	InProgress    []InProgress       `bson:"inProgress"`
	CreatedAt     time.Time          `bson:"createdAt"`
}

type Words struct {
	ID       primitive.ObjectID `bson:"_id"`
	Eng      string             `bson:"eng"`
	Pos      string             `bson:"pos"`
	Rus      string             `bson:"rus"`
	Ukr      string             `bson:"ukr"`
	Pln      string             `bson:"pln"`
	Level    string             `bson:"level"`
	Examples []Sentences        `bson:"examples"`
}

type Sentences struct {
	Ru string `bson:"ru"`
	En string `bson:"en"`
}

type InProgress struct {
	ID      primitive.ObjectID `bson:"_id"`
	WordsID primitive.ObjectID `bson:"words_id"`
}

type StudyResult struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"userID"`
	WordsID   primitive.ObjectID `bson:"wordsPairID"`
	Status    string             `bson:"status"`
	LearnedAt time.Time          `bson:"learnedAt"`
	Learned   bool               `bson:"learned"`
	Showed    bool               `bson:"shown"`
}

type Statistics struct {
	ID                  primitive.ObjectID `bson:"_id"`
	CountOfLearnedWords int                `bson:"countOfLearnedWords"`
}

type (
	GptBody struct {
		Model    string     `json:"model"`
		Stream   bool       `json:"stream"`
		Messages [1]Message `json:"messages"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	PostResponse struct {
		Choices []Choice `json:"choices"`
	}

	Choice struct {
		Index        int         `json:"index"`
		Message      messageResp `json:"message"`
		FinishReason string      `json:"finish_reason"`
	}

	messageResp struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)
