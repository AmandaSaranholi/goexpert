FROM golang:1.21.6-alpine

WORKDIR /app

COPY go.mod /
RUN go mod download

COPY cmd/ ./cmd/

RUN go build -o /stress-test ./cmd/main.go

ENTRYPOINT ["/stress-test"]