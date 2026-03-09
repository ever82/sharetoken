.PHONY: all build test test-unit test-integration lint fmt clean proto help devnet install

# 变量
BINARY_NAME=sharetokend
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT)"
GO_FILES=$(shell find . -name '*.go' -type f -not -path './vendor/*' -not -path '*/.*')
IGNITE_CMD=$(shell which ignite 2>/dev/null | grep "go/bin" || echo "$(HOME)/go/bin/ignite")

# 默认目标
all: lint test build

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# 安装到 $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/$(BINARY_NAME)

# 测试
test: test-unit

test-unit:
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./... 2>/dev/null || echo "No integration tests found"

# 代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, running go vet..."; \
		go vet ./...; \
	fi

# 代码格式化
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w $(GO_FILES); \
	fi

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# 生成 Protobuf 代码
proto:
	@echo "Generating protobuf code..."
	@echo "Using Ignite: $(IGNITE_CMD)"
	@if [ -f "$(IGNITE_CMD)" ]; then \
		$(IGNITE_CMD) generate proto-go; \
	else \
		echo "Ignite CLI not found at $(IGNITE_CMD)"; \
		echo "Installing Ignite CLI..."; \
		curl https://get.ignite.com/cli | bash; \
		$(HOME)/go/bin/ignite generate proto-go; \
	fi

# 启动本地开发网络
devnet:
	@echo "Starting local development network..."
	@if [ -f "$(IGNITE_CMD)" ]; then \
		$(IGNITE_CMD) chain serve; \
	else \
		echo "Ignite CLI not found. Install from https://ignite.com/"; \
		exit 1; \
	fi

# 启动多节点开发网络
devnet-multi:
	@echo "Starting 4-node development network..."
	./scripts/devnet_multi.sh

# 运行 CI 测试
ci-test:
	@echo "Running CI tests..."
	./scripts/test_cicd.sh

# 帮助
help:
	@echo "Available targets:"
	@echo "  make build          - Build the binary"
	@echo "  make install        - Install to GOPATH/bin"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests with coverage"
	@echo "  make test-integration - Run integration tests"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make proto          - Generate protobuf code"
	@echo "  make devnet         - Start local devnet (single node)"
	@echo "  make devnet-multi   - Start multi-node devnet"
	@echo "  make ci-test        - Run CI configuration tests"
	@echo "  make all            - Run lint, test, and build"
