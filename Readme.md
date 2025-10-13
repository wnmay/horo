# Start project

```bash
docker-compose up -d
go mod tidy
cd services/[service name]
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*$' -or $_ -match '^#') { return }
    $pair = $_ -split '=', 2
    [System.Environment]::SetEnvironmentVariable($pair[0], $pair[1])
}
go run ./cmd/main.go
```

## Initialize new service

```bash
make SERVICE=service-name
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
