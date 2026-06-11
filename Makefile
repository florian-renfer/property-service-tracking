APP_NAME=property-service-tracking
GOBIN?=$(shell go env GOPATH)/bin
REVIVE_VERSION?=v1.10.0
GOFUMPT_VERSION?=v0.7.0

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

test:
	go test -v ./...

test-race:
	go test -race ./...

test-cover:
	go test -cover ./...

lint:
	$(GOBIN)/revive -set_exit_status ./...

fmt:
	$(GOBIN)/gofumpt -w .

fmt-check:
	@out="$$( $(GOBIN)/gofumpt -l . )"; \
	if [ -n "$$out" ]; then \
		echo "$$out"; \
		exit 1; \
	fi

vet:
	go vet ./...

sqlc:
	sqlc generate

migrate-up:
	migrate -path sql/migrations \
	-database "$(DATABASE_URL)" up

migrate-down:
	migrate -path sql/migrations \
	-database "$(DATABASE_URL)" down 1

docker-up:
	docker compose up -d

docker-down:
	docker compose down

dev:
	air

tools:
	go install github.com/mgechev/revive@$(REVIVE_VERSION)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
