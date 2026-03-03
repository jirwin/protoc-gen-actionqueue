.PHONY: build generate example test test-golden clean

build:
	mkdir -p build
	go build -o build/protoc-gen-actionqueue .

generate:
	buf generate protos

example: build
	buf generate --template buf.example.gen.yaml example

test: test-golden
	go test ./...

test-golden: example
	@echo "Checking generated files match checked-in versions..."
	@git diff --exit-code example/ || \
		(echo "ERROR: Generated files differ from checked-in files. Run 'make example' and commit." && exit 1)

clean:
	rm -rf build
