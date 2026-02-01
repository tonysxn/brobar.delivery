FROM golang:1.24-alpine AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/payment-service ./payment-service/cmd/main.go

FROM alpine:3.20

RUN apk add --no-cache bash curl

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build /app/bin/payment-service .

RUN chown -R appuser:appgroup /app \
    && chmod +x /app/payment-service

USER appuser

EXPOSE ${PAYMENT_SERVICE_PORT:-8081}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:${PAYMENT_SERVICE_PORT:-8081}/health || exit 1

CMD ["/app/payment-service"]
