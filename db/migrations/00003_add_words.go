package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"io"
	parser "lWords/db/sqlc"
	"log"
	"os"
	"sync"
)

func init() {
	goose.AddMigrationNoTxContext(upAddWords, nil)
}

type Word struct {
	Eng, Rus string
}

type Words []Word

func CreateDBPool() (*pgxpool.Pool, error) {
	url := "postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords"
	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func upAddWords(ctx context.Context, tx *sql.DB) error {
	words, err := getWordsFromFile()
	if err != nil {
		log.Fatal(err)
	}

	//conn, err := pgx.Connect(ctx, "postgres://lWordsAdmin:supersecret@server:5555/lWords")
	conn, err := CreateDBPool()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer conn.Close()

	queries := parser.New(conn)

	var wg sync.WaitGroup

	for _, word := range words { // need to separate this loop, for speed up program works
		wg.Add(1)
		go func(word Word) {
			defer wg.Done()
			_, err := queries.CreateWord(context.Background(), parser.CreateWordParams{
				Eng: word.Eng, // in sqlc if field not "not null" there will be pgtype.Type{}
				Rus: word.Rus,
			})
			if err != nil {
				fmt.Println("Insert error:", err)
			}
			//fmt.Println(insertedWord)
		}(word)
	}
	wg.Wait()

	return nil
}

func getWordsFromFile() (Words, error) {
	jsonFilePath := "cmd/parser/words.json"

	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("Error read file:", err)
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
		log.Fatal("Error ReadAll:", err)
		return nil, err
	}

	var words Words
	if err = json.Unmarshal(data, &words); err != nil {
		log.Fatal("Unmarshal:", err)
		return nil, err
	}

	return words, err
}

//func downAddWords(tx *sql.Tx) error {
//	// This code is executed when the migration is rolled back.
//	return nil
//}
