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

.PHONY: tests-make-coverage
tests-make-coverage:
	go test -coverprofile=coverage.out ./...
	grep -vE "mock|\.gen\.go" coverage.out > coverage.filtered.out

.PHONY: tests-cover
tests-cover: tests-make-coverage
	go tool cover -func=coverage.filtered.out

.PHONY: tests-cover-total
tests-cover-total: tests-make-coverage
	go tool cover -func=coverage.filtered.out | grep total

.PHONY: tests-cover-html
tests-cover-html: tests-make-coverage
	go tool cover -html=coverage.filtered.out

.PHONY: clean-bin
clean-bin:
	rm -f cmd/shortener/shortener

.PHONY: create-migration
create-migration:
	# example: make create-migration name=add_short_link_table
	migrate create -ext sql -dir ./migrations -format "20060102150405" $(name)

.PHONY: mocks
mocks:
	mockery

.PHONY: emulate
emulate:
	# examples:
	#	make emulate host="http://localhost:8080"
	"./scripts/emulate-hey-requests.sh" $(host)

.PHONY: profile
profile:
	# examples:
	# 	make profile host=http://localhost:8080 out=./profiles/base.pprof
	# 	make profile host=http://localhost:8080 out=./profiles/result.pprof
	curl -sK -v $(host)/debug/pprof/profile?seconds=30 > $(out)

.PHONY: heap
heap:
	# examples:
	# 	make heap host=http://localhost:8080 out=./profiles/base.heap.pprof
	# 	make heap host=http://localhost:8080 out=./profiles/result.heap.pprof
	curl -sK -v $(host)/debug/pprof/heap?seconds=30 > $(out)

.PHONY: analyze
analyze:
	# examples:
	#	make analyze file=./profiles/base.pprof
	#	make analyze file=./profiles/result.pprof
	go tool pprof $(file)

.PHONY: analyze-web
analyze-web:
	# examples:
	#	make analyze-web file=./profiles/base.pprof
	#	make analyze-web file=./profiles/result.pprof
	go tool pprof -http=":9090" $(file)

.PHONY: analyze-diff
analyze-diff:
	# examples:
	#	make analyze-diff base=./profiles/base.pprof result=./profiles/result.pprof
	go tool pprof -top -diff_base=$(base) $(result)
