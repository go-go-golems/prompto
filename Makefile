.PHONY: all test build lint lintmax docker-lint gosec govulncheck goreleaser tag-major tag-minor tag-patch release install

all: test build

VERSION=v0.1.14

TAPES=$(shell ls doc/vhs/*tape)
gifs: $(TAPES)
	for i in $(TAPES); do vhs < $$i; done

docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run -v

lint:
	golangci-lint run -v

lintmax:
	golangci-lint run -v --max-same-issues=100

gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude=G101,G304,G301,G306,G204 -exclude-dir=.history ./...

govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	go test ./...

build:
	go generate ./...
	go build ./...

goreleaser:
	goreleaser release --skip=sign --snapshot --clean

tag-major:
	git tag -a $$(svu major) -m "Release $$(svu major)"
	git push origin $$(svu major)

tag-minor:
	git tag -a $$(svu minor) -m "Release $$(svu minor)"
	git push origin $$(svu minor)

tag-patch:
	git tag -a $$(svu patch) -m "Release $$(svu patch)"
	git push origin $$(svu patch)

release:
	git push --tags
	GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/prompto@$(shell svu current)

bump-glazed:
	go get github.com/go-go-golems/glazed@latest
	go get github.com/go-go-golems/clay@latest
	go mod tidy

BINARY=$(shell which prompto)
install:
	go build -o ./dist/prompto ./cmd/prompto && \
		cp ./dist/prompto $(BINARY)

# CodeQL local analysis target
codeql-local:
	@if [ -z "$(shell which codeql)" ]; then echo "CodeQL CLI not found. Install from https://github.com/github/codeql-cli-binaries/releases"; exit 1; fi
	@if [ ! -d "$(HOME)/codeql-go" ]; then echo "CodeQL queries not found. Clone from https://github.com/github/codeql-go"; exit 1; fi
	codeql database create --language=go --source-root=. ./codeql-db
	codeql database analyze ./codeql-db $(HOME)/codeql-go/ql/src/go/Security --format=sarif-latest --output=codeql-results.sarif
	@echo "Results saved to codeql-results.sarif"
