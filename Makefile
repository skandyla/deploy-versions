all: lint test
.DEFAULT_GOAL := all

lint:
	golangci-lint run

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.42.1

get_deps:
	go mod download

#test:
#	go test -v ./...
#	go test -race ./...
#
#cover:
#	go test -v ./... -coverprofile=/tmp/c.out
#	go tool cover -func=/tmp/c.out
