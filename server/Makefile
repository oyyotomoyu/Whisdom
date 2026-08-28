BINARY := whisdom-server

.PHONY: install run build test lint clean

install:
	go mod download

run:
	go run .

build:
	CGO_ENABLED=0 go build -o bin/$(BINARY) .

test:
	go test ./...

lint:
	@fmtout="$$(gofmt -l .)"; \
	if [ -n "$$fmtout" ]; then \
		echo "gofmt needs to be run on:"; echo "$$fmtout"; exit 1; \
	fi
	go vet ./...

clean:
	rm -rf bin
