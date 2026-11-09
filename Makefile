.PHONY: test build examples clean

test:
	go test ./...

build:
	go build -o bin/pseugo ./cmd/pseugo

examples:
	mkdir -p build
	go run ./cmd/pseugo convert -to-golang -i examples/div_mod.pseudo -o build/div_mod.go
	go run ./build/div_mod.go

clean:
	rm -rf bin build
