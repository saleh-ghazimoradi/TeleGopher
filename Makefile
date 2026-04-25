docker-up:
	docker compose up -d

docker-down:
	docker compose down

migrate-up:
	go run . migrateUp

migrate-down:
	go run . migrateDown

fmt:
	go fmt ./...

vet:
	go vet ./...

run: fmt vet
	go run . run

generate_doc:
	mkdir -p docs
	swag init -g main.go -o docs --parseDependency --parseInternal --exclude .git,docker-compose.yml,infra