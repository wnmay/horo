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


# Makefile Documentation


**Make for Windows**
   - Install via [Chocolatey](https://chocolatey.org/): `choco install make`
   - Or download [GnuWin32 Make](http://gnuwin32.sourceforge.net/packages/make.htm)

## Available Commands

### 1. Create a New Service


```powershell
# Create service with custom name
make create SERVICE=order-service
```

### 2. List Proto Files in /proto directory


```powershell
make list-protos
```

### 3. Generate Code from a Single Proto File

```powershell
# Generate from a proto file
make proto FILE=user_management.proto
```

- Outputs to `shared/proto/{package_name}/`

### 4. Generate Code from All Proto Files

```powershell
make proto-all
```

- Generates code for each file in its own package directory

---

### 5. Clean Generated Proto Files

Removes all generated protobuf files.

```powershell
# Using proto-clean
make proto-clean

# Or using clean (same effect)
make clean
```
---

### 6. Run All Services

Starts all microservices concurrently in separate processes.

```powershell
make run
```

---


### Chain Commands

```powershell
# Clean and regenerate all protos
make clean && make proto-all
```