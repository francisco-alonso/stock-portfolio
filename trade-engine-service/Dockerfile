FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o trade-engine-service ./cmd/main.go

CMD ["./trade-engine-service"]
