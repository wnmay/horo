# Start project

```bash
docker-compose up -d
go mod tidy
go run ./services/{service-name}/cmd/main.go
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

# Local Deployment with Tilt

## Prerequisites

### 1.Setup Environment by Platform

#### **Windows (WSL2 Required)**

Tilt builds code for Linux, so Windows users must use WSL2:
WSL Installation guide: https://medium.com/@sidsamanta/installing-wsl-in-windows-10-b6e8d04f5481

**Configure Docker Desktop for WSL2:**

- Open Docker Desktop Settings
- Go to Resources → WSL Integration
- Enable integration with your WSL2 distro (Ubuntu)
- Click "Apply & Restart"

#### **Window with no WSL2 (not recommended)**

Alternatively, if you don't want to use WSL, you must add .bat file to build compatible code for linux system
Example, chat-build.bat

```bash
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o build/chat-service ./services/chat-service/cmd/main.go
```

And add tilt script to run .bat file

```bash
chat_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/chat-service ./services/chat-service'
if os.name == 'nt':
  chat_compile_cmd = './infra/development/docker/chat-build.bat'
```

### 2. Install tools

#### **For Windows WSL2**

Open WSL2 terminal (Ubuntu) and run:

```bash
# Install Minikube
# See: https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Farm64%2Fstable%2Fbinary+download if this command is bug
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install Tilt
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

# Install Make
sudo apt-get update
sudo apt-get install make
```

#### **For macOS using Home Brew**

```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install Minikube
brew install minikube

# Install kubectl
brew install kubectl

# Install Tilt
brew install tilt

# Install Make (usually pre-installed, verify with: make --version)
```

## Setup Configuration

**Create centralized secrets file:**

All secrets are now centralized in a single file:

```
infra/development/k8s/secrets.yaml
```

> **Note:** This file is gitignored. Make sure to create it locally and never commit secrets to version control.

## Start Local Cluster

Run on terminal

```bash
# Start Minikube
minikube start

# Verify cluster is running
kubectl cluster-info

# Check nodes
kubectl get nodes
```

## Deploy with Tilt

```bash
tilt up

# Open Tilt UI in browser
# Visit http://localhost:10350
```

## Stop Deployment

```bash
# Stop Tilt (Ctrl+C in terminal, then)
tilt down
```

## Clean all deployment

```bash
minikube stop # Stop Minikube (optional)
minikube delete # Delete Minikube cluster (if needed)
```

# Developing with Makefile

**Make Installation**

- **Windows:** Install in WSL2: `sudo apt-get install make`
- **macOS:** Pre-installed (verify with `make --version`)

## Available Commands

**Windows WSL2:** Run in WSL2 terminal  
**macOS:** Run in Terminal

### 1. Create a New Service

```bash
# Create service with custom name
make create SERVICE=order-service
```

### 2. List Proto Files in /proto directory

```bash
make list-protos
```

### 3. Generate Code from a Single Proto File

```bash
# Generate from a proto file
make proto FILE=user_management.proto
```

- Outputs to `shared/proto/{package_name}/`

### 4. Generate Code from All Proto Files

```bash
make proto-all
```

- Generates code for each file in its own package directory

---

### 5. Clean Generated Proto Files

Removes all generated protobuf files.

```bash
# Using proto-clean
make proto-clean

# Or using clean (same effect)
make clean
```

---

### 6. Run All Services

Starts all microservices concurrently in separate processes.

```bash
make run
```

---

### Chain Commands

```bash
# Clean and regenerate all protos
make clean && make proto-all
```
