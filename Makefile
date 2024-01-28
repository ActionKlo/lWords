#migrations
mongoUp:
	go run ./cmd/mongoMigration/main.go

run:
	go run ./cmd/app/