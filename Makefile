.PHONY: proto-format proto-lint proto-gen license format lint test-unit build local-image test-e2e
all: proto-all format lint test-unit build local-image test-e2e

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building simd..."
	@cd simapp && make build 1> /dev/null
	@echo "✅ Completed build!"

###############################################################################
###                                 Tooling                                 ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

FILES := $(shell find . -name "*.go" -not -path "./simapp/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go")
license:
	@go-license --config .github/license.yml $(FILES)

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed linting!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.50
BUILDER_VERSION=0.15.3

proto-all: proto-format proto-lint proto-breaking proto-gen

proto-breaking:
	@echo "🤖 Running protobuf breaking checks..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) breaking --against "https://github.com/noble-assets/dollar.git#branch=v1.0.2"
	@echo "✅ Completed protobuf breaking checks!"

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

###############################################################################
###                                 Testing                                 ###
###############################################################################

local-image:
ifeq (,$(shell which heighliner))
	@echo heighliner not found. https://github.com/strangelove-ventures/heighliner
else
	@echo "🤖 Building image..."
	@heighliner build --chain noble-dollar-simd --file e2e/chains.yaml --local
	@echo "✅ Completed build!"
endif

test-unit:
	@echo "🤖 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -race -v . ./keeper ./types/portal/ntt
	@go tool cover -html coverage.out -o coverage.html
	@echo "✅ Completed unit tests!"

test-e2e:
	@echo "🤖 Running e2e tests..."
	@cd e2e && go test -timeout 15m -race -v ./...
	@echo "✅ Completed e2e tests!"
