# Clean Arch
## Incluido servico para listagem de ordem.

## Endpoint REST (:8000)
Executar a request do arquivo api/list_order.http

## Service GRPC (:50051)
- package pb
- service OrderService
- call ListOrders
- 
## Query GraphQL (:8080)
```sh
-query {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```