# Makefile
.PHONY: build

BINARY_NAME=projectreshoot

build:
	go mod tidy && \
   	templ generate && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}

dev:
	templ generate --watch &\
	air &\
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

clean:
	go clean
