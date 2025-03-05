# Makefile
.PHONY: build
.PHONY: migrate

BINARY_NAME=projectreshoot

build:
	tailwindcss -i ./pkg/embedfs/files/css/input.css -o ./pkg/embedfs/files/css/output.css && \
	go mod tidy && \
   	templ generate && \
	go generate ./cmd/projectreshoot && \
	go build -ldflags="-w -s" -o ./bin/${BINARY_NAME}${SUFFIX} ./cmd/projectreshoot

dev:
	templ generate --watch &\
	air &\
	tailwindcss -i ./pkg/embedfs/files/css/input.css -o ./pkg/embedfs/files/css/output.css --watch

tester:
	go mod tidy && \
	go run . --port 3232 --tester --loglevel trace

test:
	go mod tidy && \
   	templ generate && \
	go generate ./cmd/projectreshoot && \
	go test ./cmd/projectreshoot
	go test ./pkg/db
	go test ./internal/middleware

clean:
	go clean

migrate:
	go mod tidy && \
	go generate ./cmd/migrate && \
	go build -ldflags="-w -s" -o ./bin/migrate${SUFFIX} ./cmd/migrate
