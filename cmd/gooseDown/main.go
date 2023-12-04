package main

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"lWords/db/migrations"
	"log"
	"path/filepath"
)

func CreateDBPool() (*pgxpool.Pool, error) {
	url := "postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords"
	dbPool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func main() {
	conn, err := CreateDBPool()
	if err != nil {
		log.Fatal(err)
	}
	db := stdlib.OpenDBFromPool(conn)

	provider, err := goose.NewProvider(database.DialectPostgres, db, migrations.Embed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== migration list ===")
	sources := provider.ListSources()
	for _, s := range sources {
		log.Printf("%-3s %-2v %v\n", s.Type, s.Version, filepath.Base(s.Path))
	}

	stats, err := provider.Status(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n=== migration status ===")
	for _, s := range stats {
		fmt.Printf("%-3s %-2v %v\n", s.Source.Type, s.Source.Version, s.State)
	}

	fmt.Println("\n=== log migration output  ===")
	results, err := provider.DownTo(context.Background(), 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== migration results  ===")
	for _, r := range results {
		fmt.Printf("%-3s %-2v done: %v\n", r.Source.Type, r.Source.Version, r.Duration)
	}
}
