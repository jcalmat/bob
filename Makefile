PACKAGE     = bob
DATE       ?= $(shell date +%FT%T%z)
VERSION    ?= $(shell git describe --tags --always)

PKG_LIST    = $(shell go list ./... | grep -v /vendor/ | grep -v /scripts/)

GO          = go
GOLINT      = golangci-lint
GODOC       = godoc
GOFMT       = gofmt

V           = 0
Q           = $(if $(filter 1,$V),,@)
M           = $(shell printf "\033[0;35m▶\033[0m")


.PHONY: all
all: vendor build

# Executables
build: ## Build go api in bin
	$(info $(M) building bob…) @
	$Q $(GO) build \
		-o bin/$(PACKAGE)_$(VERSION)
	$Q cp bin/$(PACKAGE)_$(VERSION) bin/$(PACKAGE)

# Vendor
.PHONY: vendor
vendor: ## Create vendor directory from go.sum
	$(info $(M) running mod vendor…) @
	$Q $(GO) mod vendor

# Tidy
.PHONY: tidy
tidy: ## Update go.sum with go.mod
	$(info $(M) running mod tidy…) @
	$Q $(GO) mod tidy

# Check
.PHONY: check
check: vendor lint

# Lint
.PHONY: lint
lint: ## Run linter check on project
	$(info $(M) running $(GOLINT)…)
	$Q $(GOLINT) run

.PHONY: fmt
fmt: ## Run gofmt on project
	$(info $(M) running $(GOFMT)…) @
	$Q $(GOFMT) ./...

.PHONY: doc
doc: ## Run godoc on project
	$(info $(M) running $(GODOC)…) @
	$Q $(GODOC) ./...

.PHONY: clean
clean: ## Clean previously built binaries
	$(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin/$(PACKAGE)_*

.PHONY: version
version: ## Print current project version
	@echo $(VERSION)
