# Estágio de construção
FROM golang:1.21.6-alpine AS builder

WORKDIR /app

# Copiando o código-fonte
COPY cmd/ ./cmd/

# Compilando a aplicação
RUN go build -o /stress-test ./cmd/main.go

# Estágio final
FROM alpine:latest

# Copiando o binário da aplicação do estágio de construção
COPY --from=builder /stress-test /stress-test

# Definindo o comando padrão ao iniciar o contêiner
ENTRYPOINT ["/stress-test"]
