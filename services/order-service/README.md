# Starting

```
// you may have to run db init again to // create order db

cd services/order-service
go run ./cmd/main.go
```

# Services

# API design

## create order (Testing for Powershell)

```
$body = @{
    claims = @{
        customerId = "abc123def456"
    }
    course_id = "456e7890-e89b-12d3-a456-426614174001"
} | ConvertTo-Json -Depth 3

Invoke-RestMethod -Uri "http://localhost:3002/api/v1/orders" -Method POST -Body $body -ContentType "application/json"
```
