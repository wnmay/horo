# ------------------------------
# Go Service Project Makefile
# ------------------------------

SERVICE ?= payment-service
SERVICE_DIR := services

# Proto configuration
PROTO_SRC_DIR := proto
PROTO_OUT_DIR := shared/proto

# Use PowerShell to check for protoc plugins on Windows
PROTOC_GEN_GO := $(shell powershell -Command "Get-Command protoc-gen-go -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source")
PROTOC_GEN_GO_GRPC := $(shell powershell -Command "Get-Command protoc-gen-go-grpc -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source")

.PHONY: all create proto clean

all: create

# ------------------------------
# Create new service
# ------------------------------
create:
	@echo "Creating Go service: $(SERVICE)"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)" mkdir "$(SERVICE_DIR)\$(SERVICE)\cmd"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\adapters\inbound" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\adapters\inbound"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\adapters\outbound" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\adapters\outbound"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\app" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\app"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\domain" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\domain"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\ports\inbound" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\ports\inbound"
	@if not exist "$(SERVICE_DIR)\$(SERVICE)\internal\ports\outbound" mkdir "$(SERVICE_DIR)\$(SERVICE)\internal\ports\outbound"
	@type nul > "$(SERVICE_DIR)\$(SERVICE)\cmd\main.go"
	@echo "Project structure for $(SERVICE) created inside $(SERVICE_DIR)/"

# ------------------------------
# Generate protobufs (with package-safe structure)
# ------------------------------

# List all available proto files
list-protos:
	@echo "Available proto files:"
	@powershell -Command "Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter *.proto -Recurse -File | ForEach-Object { $$_.Name }"

# Generate proto code for a specific file
# Usage: make proto FILE=user_management.proto
# or:    make proto FILE=proto/user_management.proto
proto:
	@powershell -Command "if ('$(FILE)' -eq '') { Write-Host 'Error: FILE parameter is required'; Write-Host ''; Write-Host 'Usage: make proto FILE=<proto-file>'; Write-Host ''; Write-Host 'Available proto files:'; Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter *.proto -Recurse -File | ForEach-Object { Write-Host \"  - $($_.Name)\" }; Write-Host ''; Write-Host 'Example: make proto FILE=user_management.proto'; exit 1 }"
	@echo 🔧 Generating Go gRPC code for: $(FILE)
	@powershell -Command "if (-not (Get-Command protoc-gen-go -ErrorAction SilentlyContinue) -or -not (Get-Command protoc-gen-go-grpc -ErrorAction SilentlyContinue)) { Write-Host 'protoc-gen-go or protoc-gen-go-grpc not found in PATH.'; Write-Host '   Please run:'; Write-Host '   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest'; Write-Host '   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'; exit 1 }"
	@powershell -Command "$$protoFile = Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter '$(FILE)' -Recurse -File | Select-Object -First 1; if (-not $$protoFile) { $$protoFile = Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter '$(Split-Path -Leaf $(FILE))' -Recurse -File | Select-Object -First 1 }; if (-not $$protoFile) { Write-Host \"Error: Proto file '$(FILE)' not found in $(PROTO_SRC_DIR)\"; Write-Host ''; Write-Host 'Available proto files:'; Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter *.proto -Recurse -File | ForEach-Object { Write-Host \"  - $($_.Name)\" }; exit 1 }; $$pkgName = [System.IO.Path]::GetFileNameWithoutExtension($$protoFile.Name); $$relativePath = $$protoFile.FullName.Substring((Get-Location).Path.Length + 1); Write-Host \"Package name: $$pkgName\"; Write-Host \"Proto file: $$relativePath\"; New-Item -ItemType Directory -Force -Path $(PROTO_OUT_DIR)\$$pkgName | Out-Null; protoc --proto_path=$(PROTO_SRC_DIR) --go_out=$(PROTO_OUT_DIR)\$$pkgName --go-grpc_out=$(PROTO_OUT_DIR)\$$pkgName --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $$relativePath; if ($$LASTEXITCODE -eq 0) { Write-Host \"Protobufs generated successfully in $(PROTO_OUT_DIR)\$$pkgName\" }"

# Generate all proto files
proto-all:
	@echo 🔧 Generating Go gRPC code for all proto files
	@powershell -Command "if (-not (Get-Command protoc-gen-go -ErrorAction SilentlyContinue) -or -not (Get-Command protoc-gen-go-grpc -ErrorAction SilentlyContinue)) { Write-Host 'protoc-gen-go or protoc-gen-go-grpc not found in PATH.'; Write-Host '   Please run:'; Write-Host '   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest'; Write-Host '   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'; exit 1 }"
	@powershell -Command "Get-ChildItem -Path $(PROTO_SRC_DIR) -Filter *.proto -Recurse -File | ForEach-Object { $$protoFile = $$_; $$pkgName = [System.IO.Path]::GetFileNameWithoutExtension($$protoFile.Name); $$relativePath = $$protoFile.FullName.Substring((Get-Location).Path.Length + 1); Write-Host ('Generating ' + $$pkgName + ' from ' + $$relativePath); New-Item -ItemType Directory -Force -Path $(PROTO_OUT_DIR)\$$pkgName | Out-Null; protoc --proto_path=$(PROTO_SRC_DIR) --go_out=$(PROTO_OUT_DIR)\$$pkgName --go-grpc_out=$(PROTO_OUT_DIR)\$$pkgName --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $$relativePath; if ($$LASTEXITCODE -ne 0) { Write-Host ('❌ Failed to generate ' + $$pkgName); exit 1 } }; Write-Host '✅ All protobufs generated successfully'"

# Clean generated files
proto-clean:
	@echo "Cleaning generated proto files..."
	@powershell -Command "if (Test-Path $(PROTO_OUT_DIR)) { Remove-Item -Recurse -Force $(PROTO_OUT_DIR) }"
	@echo "Clean complete"

.PHONY: proto proto-all proto-clean list-protos

# ------------------------------
# Clean generated protos
# ------------------------------
clean:
	@echo "Cleaning generated protobuf files..."
	@powershell -Command "if (Test-Path $(PROTO_OUT_DIR)) { Remove-Item -Recurse -Force $(PROTO_OUT_DIR) }"
	@echo "Clean complete."

PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

run:
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run ./services/api-gateway/cmd/main.go'"
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run ./services/user-management-service/cmd/main.go'"
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run ./services/order-service/cmd/main.go'"
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run ./services/payment-service/cmd/main.go'"
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'go run ./services/chat-service/cmd/main.go'"
	@echo "All services started in separate windows"
run-chat:
	@powershell -Command "Start-Process powershell -NoNewWindow -ArgumentList '-NoExit', '-Command', 'go run ./services/chat-service/cmd/main.go'"

run-gateway:
	@powershell -Command "Start-Process powershell -NoNewWindow -ArgumentList '-NoExit', '-Command', 'go run ./services/api-gateway/cmd/main.go'"

run-user-management:
	@powershell -Command "Start-Process powershell -NoNewWindow -ArgumentList '-NoExit', '-Command', 'go run ./services/user-management-service/cmd/main.go'"