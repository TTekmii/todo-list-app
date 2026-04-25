.PHONY: run start-db

start-db:
	docker-compose up -d db

run:
	go run cmd/api/main.go