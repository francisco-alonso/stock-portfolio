FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o user-service ./cmd/main.go

EXPOSE 8080

CMD ["./user-service"]
