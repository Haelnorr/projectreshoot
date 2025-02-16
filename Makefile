# Makefile
.PHONY: build

BINARY_NAME=projectreshoot

build:
	go mod tidy && \
   	templ generate && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}${SUFFIX}

dev:
	templ generate --watch &\
	air &\
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

tester:
	go mod tidy && \
	go run . --port 3232 --test --loglevel trace

test:
	rm -f **/.projectreshoot-test-database.db && \
	go mod tidy && \
   	templ generate && \
	go generate && \
	go test ./middleware

clean:
	go clean
