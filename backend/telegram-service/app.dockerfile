FROM golang:1.24-alpine AS build

RUN apk update && \
    apk --no-cache add gcc g++ make git

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build the telegram-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/telegram-service ./telegram-service/cmd/telegram/main.go

# ... build stage remains same ...

FROM alpine:3.20

RUN apk update && \
    apk add --no-cache curl ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the built binary
COPY --from=build /app/bin/telegram-service .
# REMOVED .env copy
# REMOVED internal-healthcheck.sh copy

# Permissions
RUN chown -R appuser:appgroup /app \
    && chmod +x /app/telegram-service

USER appuser

EXPOSE ${TELEGRAM_PORT:-3010}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:${TELEGRAM_PORT:-3010}/health || exit 1

CMD ["/app/telegram-service"]
