# Go Sistema de Stress Test CLI

O projeto implementa um sistema para realizar testes de carga em um serviço web, gerando um relatório com as informações obtidas.


## Estrutura

```
stress-test/
├── Dockerfile
├── go.mod
└── cmd/
    └── main.go
```

## Funcionalidades

- Realizar requests HTTP para a URL especificada.
- Distribuir os requests de acordo com o nível de concorrência definido.
- Garantir que o número total de requests seja cumprido.
- Apresentar um relatório ao final dos testes contendo:
  - Tempo total gasto na execução.
  - Quantidade total de requests realizados.
  - Quantidade de requests com status HTTP 200.
  - Distribuição de outros códigos de status HTTP (como 404, 500, etc.).


## Pré-requisito

- Docker 

## Como executar

### 1. Clonar o Repositório

```sh
   git clone git@github.com:AmandaSaranholi/goexpert.git
   cd stress-test
```

### 2. Construir a Imagem Docker

```sh
docker build -t stress-test .
```

### 3 Executar o Teste de Carga

Execute o container Docker com os parâmetros desejados. Por exemplo:

```sh
docker run stress-test --url=http://google.com --requests=10 --concurrency=10
```

### Parâmetros

- `--url`: URL do serviço a ser testado.
- `--requests`: Número total de requests.
- `--concurrency`: Número de chamadas simultâneas.
- `--logs`: true/false


## Resultados da execução na maquina local

Conforme solicitado segue os resultados da execução.

Testei também com outras Urls e estou recebendo os status corretamente por aqui.

![stress-test](https://github.com/user-attachments/assets/1eecefe7-ba13-40c5-a3aa-dc7517c02267)






