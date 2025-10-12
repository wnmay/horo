FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /bin/payment-svc ./cmd/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=build /bin/payment-svc /app/payment-svc

EXPOSE 3001
ENV REST_PORT=3001

CMD ["/app/payment-svc"]
