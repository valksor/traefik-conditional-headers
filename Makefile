.PHONY: quality tools lint test coverage coverage-html yaegi_test vendor clean
.NOTPARALLEL:   # formatters mutate the tree; never race them under `make -j`

export GO111MODULE=on

default: quality test coverage

quality:
	@if command -v goimports >/dev/null; then find . -name '*.go' -not -path './.claude/*' -not -path './vendor/*' -exec goimports -w {} +; else echo "goimports not found — run 'make tools' (skipping)"; fi
	@if command -v gofumpt >/dev/null; then find . -name '*.go' -not -path './.claude/*' -not -path './vendor/*' -exec gofumpt -l -w {} +; else echo "gofumpt not found — run 'make tools' (skipping)"; fi
	go vet ./...
	@command -v gawk >/dev/null || { echo "gawk is required for the import-alias check (brew install gawk / apt-get install -y gawk)"; exit 1; }
	@alias_out="$$(./.github/alias.sh)"; rc=$$?; \
	if [ $$rc -ne 0 ]; then echo "alias check failed to run (exit $$rc)"; exit 1; fi; \
	if [ -n "$$alias_out" ]; then echo "Unnecessary import alias detected:"; echo "$$alias_out"; exit 1; fi
	@command -v golangci-lint >/dev/null || { echo "golangci-lint v2.12.2 required — see README Development Setup"; exit 1; }
	golangci-lint run ./... --fix

# Installs the go-installable tools only. golangci-lint (v2.12.2) + gawk are installed
# separately — see README Development Setup.
tools:
	go install golang.org/x/tools/cmd/goimports@v0.47.0
	go install mvdan.cc/gofumpt@v0.10.0

lint:
	golangci-lint run

test:
	go test -v -cover ./...

coverage:
	go test -race -covermode atomic -coverprofile=covprofile ./...

coverage-html: coverage
	go tool cover -html=covprofile -o coverage.html

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor

clean:
	rm -rf ./vendor
