FROM golang:1.24-alpine AS build

ENV GOPROXY=https://proxy.golang.org,direct

RUN apk --no-cache add git bash curl

WORKDIR /app

RUN go install github.com/air-verse/air@v1.61.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM golang:1.24-alpine

ENV GOPROXY=https://proxy.golang.org,direct

RUN apk --no-cache add git bash curl

WORKDIR /app

COPY --from=build /go/bin/air /usr/local/bin/air
COPY --from=build /app /app

ENV PATH="/usr/local/bin:${PATH}"

CMD ["sh", "./scripts/dev-air-runner.sh", "ontaxi-service", "cmd/ontaxi/main.go"]
