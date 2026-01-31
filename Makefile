docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	go run . migrateUp

migrate-down:
	go run . migrateDown

migrate-rollback:
	go run . migrateRollback

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run . run

