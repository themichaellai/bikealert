export GO111MODULE=on


src = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

bin/bikealert: $(src)
	go build -o bin/bikealert ./cmd/bikealert

.PHONY: lint
lint:
	golint `go list ./...`
