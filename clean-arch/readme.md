# Clean Arch

Incluido servico para listagem de ordem

# Criar a migration

docker-compose up

migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up


## Endpoint REST (:8000)

Executar a request do arquivo api/list_order.http


## Service GRPC (:50051)

- package pb
- service OrderService
- call ListOrders


## Query GraphQL (:8080)

```sh
query {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```