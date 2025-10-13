# Start project

```bash
docker-compose up -d
go mod tidy
cd services/[service name]
go run ./cmd/main.go
```

# Running each service

```bash
cd services/[service name]
go run ./cmd/main.go
```

## Initialize new service

```bash
make SERVICE=service-name
```

## Read config

from ./shared/config

```go
_ = config.LoadEnv("payment-service")
```

## Get .env

from ./shared/env

```go
port := env.GetString("REST_PORT", "3001")
num  := env.GetInt("SECRET_INT", 3000)
```
