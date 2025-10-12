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
