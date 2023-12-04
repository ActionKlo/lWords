REMOTE_SERVER = root@server
DOCKER_IMAGE = postgres:16-alpine

deploy:
	ssh $(REMOTE_SERVER) "docker run --name lWords -e POSTGRES_USER=lWordsAdmin -e POSTGRES_PASSWORD=supersecret -p 5455:5432 -d postgres"

#migrations
gooseUp:
	#goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords up
	go run ./cmd/goose # add up/down commands
gooseDown:
	#goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords down
	go run ./cmd/gooseDown
gooseStatus:
	goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords status
gooseValidate:
	goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords validate
gooseFix:
	goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords fix
gooseReset:
	goose -dir db/migrations postgres postgresql://lWordsAdmin:supersecret@100.66.158.79:5555/lWords reset

#docker
dockerComposeUp:
	scp ./docker-compose.yaml root@server:/home/GolangProjects/lWords/
	ssh $(REMOTE_SERVER) "cd /home/GolangProjects/lWords/ && docker compose up -d"
dockerComposeDown:
	ssh $(REMOTE_SERVER) "cd /home/GolangProjects/lWords/ && docker compose down"