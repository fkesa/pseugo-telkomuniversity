.PHONY: test build examples clean

test:
	go test ./...

build:
	go build -o bin/pseudoc ./cmd/pseudoc

examples:
	mkdir -p build
	go run ./cmd/pseudoc -in examples/div_mod.pseudo -out build/div_mod.go
	go run ./build/div_mod.go

clean:
	rm -rf bin build
