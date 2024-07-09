
# Go Rate Limiter

Este projeto implementa um rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso. A aplicação utiliza Redis para armazenar as informações de limitação e permite configuração via variáveis de ambiente ou arquivo .env.

## Estrutura do projeto

```
├── cmd
│   └── server
│       └── main.go
├── config
│   └── config.go
├── internal
│   └── middleware
│       └── rate_limiter.go
│   └── service
│       └── rate_limiter_service.go
├── pkg
│   └── ratelimiter
│       ├── redis.go
│       └── strategy.go
├── test
│   └── rate_limiter_test.go
├── docker-compose.yml
├── Dockerfile
├── .env
└── README.md
```

### Requisitos

- Docker
- Docker Compose
- Go 1.21.6


### Funcionalidades

- Limitação de requisições por IP 
- Limitação de requisições por token 
- Opção para configurar o tempo de bloqueio após exceder o limite 
- Armazenamento das informações de limitação no Redis.
- Middleware para fácil integração com servidores web.
- Suporte para troca de mecanismo de persistência.


### Configuração do Ambiente

1. **Clone o repositório**

   ```sh
   git clone git@github.com:AmandaSaranholi/goexpert.git
   cd rate-limiter
   ```

2. **Crie e configure o aquivo .env**

   ```sh
   cp .env.example .env
   # Modifique o arquivo .env como desejar
   ```

3. **Inicie o servidor Redis e o servidor Go com Docker Compose**

   ```sh
   docker-compose up -d
   ```

4. **O servidor vai responder no endereço `http://localhost:8080`**

5. **Faça requisições http para a URL do servidor usando umo navegador, HTTP Client ou até mesmo cURL**

   Exemplo cURL com Header

   ```sh
   curl --request GET \
   --url http://localhost:8080/ \
   --header 'API_KEY: xxx123'
   ```

### Rodando os testes

1. **Execute os testes**
   ```sh
   docker compose exec app go test -v ./test
   ```

## Configuração

Configuracao do arquivo `.env`:

```
STRATEGY_TYPE=redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
LIMITER_IP_LIMIT=5
LIMITER_TOKEN_LIMIT=5
LIMITER_BLOCK_DURATION=300

```

## Padrão Strategy

O rate limiter utiliza o padrão de projeto Strategy, permitindo a fácil troca entre diferentes mecanismos de persistência. 

## Middleware

O rate limiter é implementado como um middleware para controlar a taxa de requisições recebidas com base no endereço IP ou token de acesso.

### Conclusão

Implementação do desafio tecnico proposto.


## Licença

[MIT License](LICENSE).