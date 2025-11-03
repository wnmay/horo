FROM alpine:3.20

WORKDIR /app

# Copy the pre-compiled binary
COPY build/order-service /app/build/order-service

# Copy shared dependencies
COPY shared /app/shared

EXPOSE 3002
ENV REST_PORT=3002

CMD ["/app/build/order-service"]