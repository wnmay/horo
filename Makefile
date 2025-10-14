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

# List all available proto files
list-protos:
	@echo "Available proto files:"
	@find $(PROTO_SRC_DIR) -name "*.proto" -type f | sed 's|$(PROTO_SRC_DIR)/||' | nl

# Generate proto code for a specific file
# Usage: make proto FILE=user_management.proto
# or:    make proto FILE=proto/user_management.proto
proto:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE parameter is required"; \
		echo ""; \
		echo "Usage: make proto FILE=<proto-file>"; \
		echo ""; \
		echo "Available proto files:"; \
		find $(PROTO_SRC_DIR) -name "*.proto" -type f | sed 's|$(PROTO_SRC_DIR)/||' | sed 's/^/  - /'; \
		echo ""; \
		echo "Example: make proto FILE=user_management.proto"; \
		exit 1; \
	fi
	@echo "🔧 Generating Go gRPC code for: $(FILE)"
	@if [ -z "$(shell which protoc-gen-go)" ] || [ -z "$(shell which protoc-gen-go-grpc)" ]; then \
		echo "protoc-gen-go or protoc-gen-go-grpc not found in PATH."; \
		echo "   Please run:"; \
		echo "   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; \
		echo "   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; \
		exit 1; \
	fi
	$(eval PROTO_FILE := $(shell find $(PROTO_SRC_DIR) -name "$(FILE)" -o -name "$(notdir $(FILE))" | head -n 1))
	@if [ -z "$(PROTO_FILE)" ]; then \
		echo "Error: Proto file '$(FILE)' not found in $(PROTO_SRC_DIR)"; \
		echo ""; \
		echo "Available proto files:"; \
		find $(PROTO_SRC_DIR) -name "*.proto" -type f | sed 's|$(PROTO_SRC_DIR)/||' | sed 's/^/  - /'; \
		exit 1; \
	fi
	$(eval PKG_NAME := $(basename $(notdir $(PROTO_FILE))))
	@echo "Package name: $(PKG_NAME)"
	@mkdir -p $(PROTO_OUT_DIR)/$(PKG_NAME)
	@protoc \
		--proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PROTO_OUT_DIR)/$(PKG_NAME) \
		--go-grpc_out=$(PROTO_OUT_DIR)/$(PKG_NAME) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)
	@echo "Protobufs generated successfully in $(PROTO_OUT_DIR)/$(PKG_NAME)"

# Generate all proto files
proto-all:
	@echo "🔧 Generating Go gRPC code for all proto files"
	@if [ -z "$(shell which protoc-gen-go)" ] || [ -z "$(shell which protoc-gen-go-grpc)" ]; then \
		echo "protoc-gen-go or protoc-gen-go-grpc not found in PATH."; \
		echo "   Please run:"; \
		echo "   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; \
		echo "   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"; \
		exit 1; \
	fi
	@for proto_file in $$(find $(PROTO_SRC_DIR) -name "*.proto"); do \
		pkg_name=$$(basename $$proto_file .proto); \
		echo "Generating $$pkg_name from $$proto_file"; \
		mkdir -p $(PROTO_OUT_DIR)/$$pkg_name; \
		protoc \
			--proto_path=$(PROTO_SRC_DIR) \
			--go_out=$(PROTO_OUT_DIR)/$$pkg_name \
			--go-grpc_out=$(PROTO_OUT_DIR)/$$pkg_name \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			$$proto_file; \
	done
	@echo "All protobufs generated successfully"

# Clean generated files
proto-clean:
	@echo "Cleaning generated proto files..."
	@rm -rf $(PROTO_OUT_DIR)
	@echo "Clean complete"

.PHONY: proto proto-all proto-clean list-protos

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
	go run ./services/order-service/cmd/main.go &
	go run ./services/payment-service/cmd/main.go &
	wait