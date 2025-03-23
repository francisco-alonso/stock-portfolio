FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o portfolio-service ./cmd/main.go

EXPOSE 8082

CMD ["./portfolio-service"]
