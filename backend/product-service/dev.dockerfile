# Stage 1: build migrate binary and install dependencies + install air
FROM golang:1.24-alpine AS build

RUN apk --no-cache add git bash curl postgresql-client

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Stage 2: dev container
FROM golang:1.24-alpine

RUN apk --no-cache add git bash curl postgresql-client

WORKDIR /app

COPY --from=build /go/bin/migrate /usr/local/bin/migrate
COPY --from=build /go/bin/air /usr/local/bin/air
COPY --from=build /app /app

RUN chmod +x /app/scripts/db-entrypoint.sh && \
    find /app -type f -name ".air.toml" -exec chmod +r {} +

ENV PATH="/usr/local/bin:${PATH}"

ENTRYPOINT ["/app/scripts/db-entrypoint.sh"]
CMD ["air", "-c", "product-service/.air.toml"]
