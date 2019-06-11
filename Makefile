SHELL := /bin/sh
GOPROCS := 4
COVFILE := coverage.out

.PHONY: get-deps
get-deps:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: cov
cov: $(COVFILE)
	go tool cover -func=$(COVFILE)
	sed -i '\|github.com/cycloidio/raws/generate.go|d' $(COVFILE)

.PHONY: htmlcov
htmlcov: $(COVFILE)
	@go tool cover -html=$(COVFILE)

$(COVFILE):
	@GO111MODULE=on go test ./... -covermode=count -coverprofile=$(COVFILE)

.PHONY: travis-ci
travis-ci: lintcheck test cov

.PHONY: test
test:
	@GO111MODULE=on go test ./... -coverprofile=$(COVFILE)

.PHONY: lintcheck
lintcheck:
	@GO111MODULE=on golangci-lint run ./... -E goimports && \
		golint -set_exit_status ./...

.PHONY: generate
generate:
	@GO111MODULE=on go generate ./...
