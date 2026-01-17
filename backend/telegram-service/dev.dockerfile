FROM golang:1.24-alpine

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk update && \
    apk --no-cache add git bash curl tdlib tdlib-dev gcompat

WORKDIR /app

ENV CGO_ENABLED=1

RUN go install github.com/air-verse/air@v1.61.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x /app/scripts/db-entrypoint.sh && \
    find /app -type f -name ".air.toml" -exec chmod +r {} +

CMD ["air", "-c", "telegram-service/.air.toml"]
