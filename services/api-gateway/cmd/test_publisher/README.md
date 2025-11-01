0. Replace rabbit uri in test with real rabbiq uri
```go
rabbitmqURI := "replace-rabbit-url-here"
```
1. Start API gateway
```bash
go run ./services/api-gateway/cmd/main.go
```
2. Start service
```bash
go run ./services/user-management-service/cmd/main.go
go run ./services/chat-service/cmd/main.go
```
3.Start test
```bash
go run ./services/api-gateway/cmd/test_publisher/main.go
```