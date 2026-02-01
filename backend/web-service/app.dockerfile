FROM golang:1.24-alpine AS build

RUN apk --no-cache add git bash postgresql-client

WORKDIR /app

# Install migration tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

RUN go build -o web-service-app ./web-service/cmd/web/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash postgresql-client

WORKDIR /app

COPY --from=build /go/bin/migrate /usr/local/bin/migrate
COPY --from=build /app/web-service-app .
COPY scripts /app/scripts
COPY web-service/migrations /app/migrations

RUN chmod +x /app/scripts/db-entrypoint.sh

ENTRYPOINT ["/app/scripts/db-entrypoint.sh"]
CMD ["./web-service-app"]
