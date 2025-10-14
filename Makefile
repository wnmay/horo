# ------------------------------
# Go Service Project Makefile
# ------------------------------

SERVICE ?= payment-service
SERVICE_DIR := services

# Proto configuration
PROTO_SRC_DIR := proto
PROTO_OUT_DIR := shared/proto

PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)

.PHONY: all create proto clean

all: create

# ------------------------------
# Create new service
# ------------------------------
create:
	@echo "Creating Go service: $(SERVICE)"
	cd $(SERVICE_DIR) && \
	mkdir -p $(SERVICE)/cmd && \
	mkdir -p $(SERVICE)/internal/adapters/inbound && \
	mkdir -p $(SERVICE)/internal/adapters/outbound && \
	mkdir -p $(SERVICE)/internal/app && \
	mkdir -p $(SERVICE)/internal/domain && \
	mkdir -p $(SERVICE)/internal/ports/inbound && \
	mkdir -p $(SERVICE)/internal/ports/outbound && \
	touch $(SERVICE)/cmd/main.go
	@echo "Project structure for $(SERVICE) created inside $(SERVICE_DIR)/"

# ------------------------------
# Generate protobufs (with package-safe structure)
# ------------------------------
# --- Proto Generation Configuration ---
PROTO_SRC_DIR := proto
PROTO_OUT_DIR := shared/proto
PROTO_PKG ?= usermanagement_v1  

PROTOC_GEN_GO := $(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC := $(shell which protoc-gen-go-grpc)

# --- Generate proto files ---
proto:
	@echo "🔧 Generating Go gRPC code for package: $(PROTO_PKG)"
	@if [ -z "$(PROTOC_GEN_GO)" ] || [ -z "$(PROTOC_GEN_GO_GRPC)" ]; then \
		echo " protoc-gen-go or protoc-gen-go-grpc not found in PATH."; \
		echo "   Please run:"; \
		echo "   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; \
		echo "   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; \
		exit 1; \
	fi
	@mkdir -p $(PROTO_OUT_DIR)/$(PROTO_PKG)
	@protoc \
		--proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PROTO_OUT_DIR)/$(PROTO_PKG) \
		--go-grpc_out=$(PROTO_OUT_DIR)/$(PROTO_PKG) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(shell find $(PROTO_SRC_DIR) -name "*.proto")
	@echo " Protobufs generated successfully in $(PROTO_OUT_DIR)/$(PROTO_PKG)"

# ------------------------------
# Clean generated protos
# ------------------------------
clean:
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(PROTO_OUT_DIR)
	@echo "Clean complete."

PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .


run:
	go run ./services/api-gateway/cmd/main.go &
	go run ./services/user-management-service/cmd/main.go &
	wait