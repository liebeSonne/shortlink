.PHONY: all
all: build tests

.PHONY: build
build:
	go build -o cmd/shortener/shortener ./cmd/shortener

.PHONY: tests
tests:
	go test ./...

.PHONY: tests-v
tests-v:
	go test -v ./...

.PHONY: clean-bin
clean-bin:
	rm -f cmd/shortener/shortener