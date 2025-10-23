# Start project

```bash
docker-compose up -d
go mod tidy
go run ./cmd/main.go
```

## Read env

from ./shared/env

```go
_ = env.LoadEnv("payment-service")
```

## Get .env

from ./shared/env

```go
port := env.GetString("REST_PORT", "3001")
num  := env.GetInt("SECRET_INT", 3000)
```

## Generate pb file

from Makefile

```bash
make proto PROTO_PKG=pkgname
```

# Makefile Documentation

## Prerequisites
**Make for Windows**
   - Install via [Chocolatey](https://chocolatey.org/): `choco install make`
   - Or download [GnuWin32 Make](http://gnuwin32.sourceforge.net/packages/make.htm)

## Available Commands

### 1. Create a New Service

**Usage:**
```powershell
# Create service with custom name
make create SERVICE=order-service
```

### 2. List Proto Files

Lists all available `.proto` files in the proto directory.

**Usage:**
```powershell
make list-protos
```

---

### 3. Generate Code from a Single Proto File

Generates Go code from a specific `.proto` file.

**Usage:**
```powershell
# Generate from a proto file
make proto FILE=user_management.proto

# Or with full path
make proto FILE=proto/user_management.proto
```

- Generates both `.pb.go` (message definitions) and `_grpc.pb.go` (gRPC service definitions) files
- Outputs to `shared/proto/{package_name}/`
---

### 4. Generate Code from All Proto Files

Generates Go code from all `.proto` files in the proto directory.

**Usage:**
```powershell
make proto-all
```

**What it does:**
- Finds all `.proto` files recursively in the proto directory
- Generates code for each file in its own package directory

---

### 5. Clean Generated Proto Files

Removes all generated protobuf files.

**Usage:**
```powershell
# Using proto-clean
make proto-clean

# Or using clean (same effect)
make clean
```

**What it does:**
- Deletes the entire `shared/proto/` directory

---

### 6. Run All Services

Starts all microservices concurrently in separate processes.

**Usage:**
```powershell
make run
```
**Note:** Services run in background processes. To stop them, you'll need to manually close the PowerShell windows or use Task Manager.

---

### Override Default Directories

You can override variables when running make commands:

```powershell
# Use custom proto directory
make proto-all PROTO_SRC_DIR=api/protos

# Use custom output directory
make proto FILE=user.proto PROTO_OUT_DIR=generated/proto

# Create service in custom location
make create SERVICE=my-service SERVICE_DIR=microservices
```

### Chain Commands

Combine multiple targets:

```powershell
# Clean and regenerate all protos
make clean && make proto-all
```