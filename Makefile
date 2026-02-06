.PHONY: build generate example test clean

build:
	mkdir -p build
	go build -o build/protoc-gen-actionqueue .

generate:
	buf generate protos

example: build
	buf generate --template buf.example.gen.yaml example

test:
	go test ./...

clean:
	rm -rf build
