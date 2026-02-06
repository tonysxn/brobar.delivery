FROM golang:1.24-alpine AS build

RUN apk update && \
    apk --no-cache add gcc g++ make git

WORKDIR /app

# Copy dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build the ontaxi-service binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/ontaxi-service ./ontaxi-service/cmd/ontaxi/main.go

FROM alpine:3.20

RUN apk update && \
    apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the built binary
COPY --from=build /app/bin/ontaxi-service .

# Permissions
RUN chown -R appuser:appgroup /app \
    && chmod +x /app/ontaxi-service

USER appuser

CMD ["/app/ontaxi-service"]
