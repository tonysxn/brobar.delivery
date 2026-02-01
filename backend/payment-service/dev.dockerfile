FROM golang:1.24-alpine

# Install git and air for live reloading
RUN apk add --no-cache git make
RUN go install github.com/air-verse/air@v1.61.0

WORKDIR /app

# Copy generic go.mod and go.sum (from parent backend dir) happens via volume mount in docker-compose
# But we need to setup the workspace

CMD ["air", "-c", "payment-service/.air.toml"]
