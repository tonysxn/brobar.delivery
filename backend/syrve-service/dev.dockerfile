FROM golang:1.24-alpine

RUN apk --no-cache add git bash curl

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x /app/scripts/db-entrypoint.sh && \
    find /app -type f -name ".air.toml" -exec chmod +r {} +

CMD ["air", "-c", "syrve-service/.air.toml"]
