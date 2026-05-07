.PHONY: build test lint fmt clean

build:
	go build ./cmd/maeve

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -s -w .

clean:
	rm -rf bin dist

