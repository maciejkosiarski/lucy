APP_NAME = lucy
BASE_GO_IMAGE = golang:1.17.6-alpine3.15
BASE_TARGET_IMAGE = alpine:3.15

REGISTRY ?= localhost:5000
IMAGE_DEV = $(REGISTRY)/$(APP_NAME)-dev
IMAGE_PROD = $(REGISTRY)/$(APP_NAME)-application


DOCKERFILE_DEV = .docker/dev/Dockerfile
DOCKERFILE_PROD = ./Dockerfile


CGO_ENABLED = 0 # statically linked = 0
TARGETOS=linux
ifeq ($(OS),Windows_NT) 
    TARGETOS := Windows
else
    TARGETOS := $(shell sh -c 'uname 2>/dev/null || echo Unknown' | tr '[:upper:]' '[:lower:]')
endif
TARGETARCH = amd64

.DEFAULT_GOAL = help
PID = /tmp/serving.pid
DEVELOPER_UID     ?= $(shell id -u)
#-----------------------------------------------------------------------------------------------------------------------
ARG := $(word 2, $(MAKECMDGOALS))
%:
	@:
#-----------------------------------------------------------------------------------------------------------------------
#-----------------------------------------------------------------------------------------------------------------------

help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

build: build-image-dev ## Alias for 'build-image-dev'

build-image-dev: ## Build dev image
	@docker build 											\
	    -t $(IMAGE_DEV)										\
		--build-arg BASE_GO_IMAGE=$(BASE_GO_IMAGE)			\
		--build-arg DEVELOPER_UID=$(DEVELOPER_UID)			\
		--build-arg APP_NAME=$(APP_NAME)					\
		-f $(DOCKERFILE_DEV)								\
		.

build-image-prod: ## Build prod image
	@docker build 											\
	    -t $(IMAGE_PROD)									\
		--build-arg BASE_GO_IMAGE=$(BASE_GO_IMAGE)			\
		--build-arg BASE_TARGET_IMAGE=$(BASE_TARGET_IMAGE)	\
		--build-arg CGO_ENABLED=${CGO_ENABLED}				\
		--build-arg TARGETOS=${TARGETOS}					\
		--build-arg TARGETARCH=${TARGETARCH}				\
		--build-arg APP_NAME=${APP_NAME}					\
		-f $(DOCKERFILE_PROD)								\
		.

run-image: ## Run prod image
	@docker run $(IMAGE_PROD)

up: ## Start application dev container
	@cd .docker && \
	COMPOSE_PROJECT_NAME=$(APP_NAME) \
	IMAGE_DEV=$(IMAGE_DEV) \
	APP_NAME=$(APP_NAME) \
	docker-compose up -d

down: ## Remove application dev container
	@cd .docker && \
	COMPOSE_PROJECT_NAME=$(APP_NAME) \
	IMAGE_DEV=$(IMAGE_DEV) \
	APP_NAME=$(APP_NAME) \
	docker-compose down

console: ## Enter application dev container
	@docker exec -it $(APP_NAME)-dev bash

go: go-build ## Alias for 'go-build'

go-build: ## Build dev application (go build)	
	@go mod tidy
	@if [ "${ARG}" = '' ]; then @env CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-X main.env=dev" -o bin/${APP_NAME} ./; fi
	@env CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-X main.env=dev" -o bin/${APP_NAME} ./

go-build-mac:
	@go mod tidy
	@env CGO_ENABLED=${CGO_ENABLED} GOOS=darwin GOARCH=${TARGETARCH} go build -ldflags "-X main.env=dev" -o bin/${APP_NAME} ./

clean: ## Clean bin/
	@rm -rf bin/${APP_NAME}
