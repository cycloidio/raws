SHELL := /bin/sh
GOPROCS := 4
SRC := $(wildcard *.go)
COVFILE := coverage.out
GOFILES_NOVENDOR := $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./example/*")

.PHONY: get-deps
get-deps:
	go get -t ./...
	go get golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: cov
cov: $(COVFILE)
	go tool cover -func=$(COVFILE)

.PHONY: htmlcov
htmlcov: $(COVFILE)
	go tool cover -html=$(COVFILE)

$(COVFILE):
	go test $(PKG) -covermode=count -coverprofile=$(COVFILE)

.PHONY: travis-ci
travis-ci: test vetcheck fmtcheck lintcheck cov

.PHONY: test
test:
	go test -v $(PKG) -coverprofile=$(COVFILE)

.PHONY: fmtcheck
fmtcheck:
	@if [ "$(shell goimports -l $(GOFILES_NOVENDOR) | wc -l)" != "0" ]; then \
		echo "Files missing go fmt: $(shell goimports -l $(GOFILES_NOVENDOR))"; exit 2; \
	else \
	    echo -e "ok\tall files passed go fmt"; \
	fi

.PHONY: vetcheck
vetcheck:
ifeq ($(shell go tool vet -all -shadow=true . 2>&1 | wc -l), 0)
	@printf "ok\tall files passed go vet\n"
else
	@echo -e "error\tsome files did not pass go vet\n"
	@go tool vet -all -shadow=true . 2>&1
endif

.PHONY: lintcheck
lintcheck:
ifeq ($(shell golint 2>&1 | wc -l), 0)
	@printf "ok\tall files passed golint\n"
else
	@echo -e "error\tsome files did not pass golint\n"
	@golint 2>&1
endif
