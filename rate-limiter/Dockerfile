FROM golang:1.21.6-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o server ./cmd/server

EXPOSE 8080

ENTRYPOINT ["./server"]