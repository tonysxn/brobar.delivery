FROM golang:1.24-alpine
ENV GOPROXY=https://proxy.golang.org,direct

RUN apk --no-cache add git bash curl postgresql-client

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@v1.61.0

# Install migration tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY scripts /app/scripts
COPY web-service/migrations /app/web-service/migrations

# Make script executable
RUN chmod +x /app/scripts/db-entrypoint.sh

ENTRYPOINT ["/app/scripts/db-entrypoint.sh"]
CMD ["sh", "./scripts/dev-air-runner.sh", "web-service", "cmd/web/main.go"]
