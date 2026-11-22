.PHONY: test build examples clean

test:
	go test ./...

build:
	go build -o bin/pseugo ./cmd/pseugo

examples:
	go run ./cmd/pseugo run examples/div_mod.pseudo

clean:
	rm -rf bin build
