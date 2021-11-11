all: lint test
.DEFAULT_GOAL := all

IMAGENAME=deploy-versions
BUILD_VERSION=$(shell git describe --tags --always)
BUILD_COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

lint:
	golangci-lint run

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.42.1

get_deps:
	go mod download

test:
	go test -v ./...
	go test -race ./...

cover:
	go test -v ./... -coverprofile=/tmp/c.out
	go tool cover -func=/tmp/c.out

build:
	@echo MARK: build go code
	GOOS=linux CGO_ENABLED=0 go build -v --ldflags "-extldflags '-static' -X app.buildVersion=$(BUILD_VERSION) -X app.buildCommit=$(BUILD_COMMIT) -X app.buildTime=$(BUILD_TIME)" -o main

docker_image:
	@echo MARK: package built code inside docker image
	docker version
	docker build -t ${IMAGENAME} --build-arg GIT_COMMIT=${BUILD_COMMIT}  -f Dockerfile .

docker_compose_start:
	@echo "MARK: Starting docker-compose"
	@docker compose up -d
	@echo "MARK: Sleeping 5 seconds until dependencies start."
	@sleep 5

docker_compose_stop:
	@echo MARK: Stopping cluster
	@docker compose down --remove-orphans

# http (httpie required)
# example. go integration tests TBD
docker_compose_run_tests:
	@echo MARK: Testing via docker-compose 
	curl -v localhost:8080/info
	curl -v localhost:8080/versions
	curl -v localhost:8080/version/130
	curl -v -X POST localhost:8080/version -d '{"project":"example","env":"dev","region":"us-west-1","service_name"="notifications","user_name":"bsmith","build_id"=133}'
	curl -v localhost:8080/version/133
	curl -v localhost:8080/version/140  #not found expected

