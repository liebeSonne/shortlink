.PHONY: all
all: build generate tests

.PHONY: build
build:
	go build -o cmd/shortener/shortener ./cmd/shortener

.PHONY: generate
generate:
	go generate ./...

.PHONY: tests
tests:
	go test ./...

.PHONY: tests-v
tests-v:
	go test -v ./...

.PHONY: clean-bin
clean-bin:
	rm -f cmd/shortener/shortener