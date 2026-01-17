FROM golang:1.24-alpine AS build

RUN apk --no-cache add gcc g++ make git postgresql-client

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/user-service ./user-service/cmd/user/main.go

# ... build stage remains same ...

FROM alpine:3.20

RUN apk add --no-cache postgresql-client bash curl

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy necessary files from build stage
COPY --from=build /app/bin/user-service .
COPY --from=build /app/user-service/migrations ./migrations
COPY --from=build /go/bin/migrate /usr/local/bin/migrate
COPY --from=build /app/scripts/db-entrypoint.sh /usr/local/bin/

# Verify the binary exists and is executable
RUN chown -R appuser:appgroup /app \
    && chmod +x /app/user-service \
    && chmod +x /usr/local/bin/db-entrypoint.sh

USER appuser

EXPOSE ${USER_PORT:-3000}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:${USER_PORT:-3000}/health || exit 1

ENTRYPOINT ["/usr/local/bin/db-entrypoint.sh"]
CMD ["/app/user-service"]