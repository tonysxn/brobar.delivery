FROM golang:1.24-alpine

RUN apk --no-cache add git bash curl postgresql-client

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@v1.61.0

# Install migration tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY scripts /app/scripts
COPY migrations /app/migrations

# Make script executable
RUN chmod +x /app/scripts/db-entrypoint.sh

ENTRYPOINT ["/app/scripts/db-entrypoint.sh"]
CMD ["air", "-c", "web-service/.air.toml"]
