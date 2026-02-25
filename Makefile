IMG_NAME=haihoanguci/user-service
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := temporary

ifeq ($(BRANCH), main)
	IMG_TAG := dev
endif

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

COVERAGE_EXCLUDE=infrastructure|mocks|vendor|test|docs|main.go|config.go|client.go
COVERAGE_THRESHOLD = 80
COVERAGE_FOLDER=./coverage
#=========================== DEV TOOLS =========================== 
.PHONY: mock-gen
mock-gen:
	go generate ./...

.PHONY: dev-up, dev-down, dev_run, swag-gen
swag-gen:
	swag init -g ./cmd/api/main.go --output ./docs

dev-up:
	docker-compose -f docker-compose.dev.yaml up -d

dev-down:
	docker-compose -f docker-compose.dev.yaml down

dev-run: swag-gen
	APP_HOST_NAME=localhost:8080 APP_PORT=:8080 DB_NAME=user go run ./cmd/api/main.go

.PHONY: test 
test: clean
	mkdir -p $(COVERAGE_FOLDER)
	go test ./... -coverprofile=$(COVERAGE_FOLDER)/coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -v -E "$(COVERAGE_EXCLUDE)" $(COVERAGE_FOLDER)/coverage.tmp > $(COVERAGE_FOLDER)/coverage.out
	go tool cover -html=$(COVERAGE_FOLDER)/coverage.out -o $(COVERAGE_FOLDER)/coverage.html
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi

.PHONY: redis-run redis-cli redis-monitor
redis-run:
	docker run --name redis -p 6379:6379 -d redis

redis-cli:
	docker exec -it redis redis-cli

redis-monitor:
	docker exec -it redis redis-cli monitor

.PHONY: docker-build, docker-up, docker-down, docker-release, docker-test
docker-build:
	@echo "Building Docker image: $(IMG_NAME):$(IMG_TAG)"
	docker build -t $(IMG_NAME):$(IMG_TAG) .

	@echo "Building Docker image for migration: $(IMG_NAME)_migration:$(IMG_TAG)"
	docker build --target migration -t $(IMG_NAME)_migration:$(IMG_TAG) .

docker-release: docker-build
	@echo "Pushing Docker image: $(IMG_NAME):$(IMG_TAG)"
	docker push $(IMG_NAME):$(IMG_TAG)

	@echo "Pushing Docker image for migration: $(IMG_NAME)_migration:$(IMG_TAG)"
	docker push $(IMG_NAME)_migration:$(IMG_TAG)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

DOCKER_HUB_USERNAME ?=
DOCKER_HUB_ACCESS_TOKEN ?=

docker-login:
	echo "$(DOCKER_HUB_ACCESS_TOKEN)" | docker login -u "$(DOCKER_HUB_USERNAME)" --password-stdin


docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark_service:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
    if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
    else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
   	fi	

.PHONY: clean
clean:
	go clean -testcache
	rm -rf $(COVERAGE_FOLDER)
# 	docker rm -f redis || true

generate-rsa-key:
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private_key.pem -out public_key.pem

#=========================== DB MIGRATION ===========================
.PHONY: new-schema
new-schema:
	migrate create -ext sql -dir ./migrations -seq $(name)
# example: make new-schema name=add_bookmark

migrate:
	go run ./cmd/migrate/main.go