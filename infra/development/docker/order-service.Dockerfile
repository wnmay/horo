FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /bin/order-svc ./services/order-service/cmd/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=build /bin/order-svc /app/order-svc

EXPOSE 3001
ENV REST_PORT=3001

CMD ["/app/order-svc"]