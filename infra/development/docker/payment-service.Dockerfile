FROM alpine

WORKDIR /app

COPY ./build/payment-service /app/build/payment-service

COPY ./shared /app/shared

EXPOSE 3001

ENTRYPOINT ["/app/build/payment-service"]

