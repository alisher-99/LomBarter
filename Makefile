.PHONY: tidy lint test swag bin-deps fmt migrate-up migrate-down migrate-down-all migrate-new compose compose-down

lint:
	golangci-lint run

tidy:
	go mod tidy
	go mod vendor

test:
	go test -v -cover -race ./internal/...

swag:
	swag init -g cmd/app/main.go -o ./swagger/ --parseVendor --exclude ./vendor

bin-deps:
	go install -tags 'cassandra' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest

fmt:
	gofumpt -w .

# Запуск миграций
migrate-up:
	migrate -path ./migrations -database "cassandra://localhost:9042/region" up

# Откатить последнюю миграцию
migrate-down:
	migrate -path ./migrations -database "cassandra://localhost:9042/region" down 1

# Откатить все миграции
migrate-down-all:
	migrate -path ./migrations -database "cassandra://localhost:9042/region" down all

# Создать новую миграцию
migrate-new:
	migrate -path ./migrations create -dir ./migrations -ext cql $(name)

compose:
	docker-compose up --build -d mongo dragonfly jaeger zookeeper kafka service servicemesh-mock-server

compose-down:
	docker-compose down --remove-orphans

integration:
	docker-compose build integration ;\
	test_status=0 ;\
	docker-compose run integration || test_status=$$? ;\
	docker-compose down --remove-orphans ;\
	echo "status="$$test_status; exit $$test_status ;\

mock:
	mockgen -source="internal/domain/repository/db_repo.go" -destination="internal/domain/repository/mock_repo/mock_db_repo.go" -package="mock_repo"
	mockgen -source="internal/domain/repository/cache_repo.go" -destination="internal/domain/repository/mock_repo/mock_cache_repo.go" -package="mock_repo"
	mockgen -source="pkg/metrics/metrics.go" -destination="pkg/metrics/mock_metrics/mock_metrics.go" -package="mock_metrics"
	mockgen -source="pkg/metrics/types/types.go" -destination="pkg/metrics/mock_metrics/mock_types.go" -package="mock_metrics"
	mockgen -source="internal/transport/broker/processor.go" -destination="internal/transport/broker/mock_broker/mock_consumer.go" -package="mock_broker"

all: mock tidy lint test swag integration