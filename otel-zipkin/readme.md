# OTEL Distributed Tracing + Zipkin

Esse projeto implementa um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin)  juntamente com a cidade. Esse sistema implementa também OTEL (Open Telemetry) e Zipkin para tracing distribuído.

## Requisitos

 - Go 1.21.6 
 - Docker 
 - Docker Compose

# Definições
 - O Serviço A será denominado `input-service`, o qual receberá o input do usuario e encaminhará para o Serviço B
 - O Serviço B será denominado  `orchestrator-service`, o qual fara a busca na API de Cep e Tempo.


### Serviço A

- Receberá um input de 8 dígitos via POST:

  ```json
  {
  	"cep": "29902555"
  }
  ```

### Serviço B 

- Receber um CEP válido de 8 dígitos.
- Realizar a pesquisa do CEP e encontrar o nome da localização.
- Retornar as temperaturas formatadas em Celsius, Fahrenheit e Kelvin, juntamente com o nome da localização.

## Implementação do OTEL + Zipkin

- Implementar tracing distribuído entre Serviço A e Serviço B.
- Utilizar spans para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura.



## Configuração do Ambiente

1. Clone o repositório:

   ```bash
   git clone git@github.com:AmandaSaranholi/goexpert.git
   cd otel-zipkin
   ```

2. Crie e configure o aquivo .env**

   ```bash
   cp .env.example .env
   # Modifique o arquivo .env como desejar
   ```

3. Inicie o servidor Go com Docker Compose:

   ```bash
   docker-compose up --build
   ```

4. O servidor vai responder no endereço: `http://localhost:8080/weather`.

5. Para testar o Serviço A, utilize um cliente HTTP (como Postman) para enviar um `POST` com o payload a seguir:

   ```json
   {
   	"cep": "17280276"
   }
   ```

6. Para visualizar os traces no `Zipkin` acesse `http://localhost:9411`.

## Resultados após os testes
![success_zipcode](https://github.com/user-attachments/assets/baa1e7b9-c25b-4a10-8965-1fd5d05554a3)
![invalid_zipcode](https://github.com/user-attachments/assets/98a7614a-9c3c-4a0b-a206-908f7c5a0878)
![canot_find_zipcode](https://github.com/user-attachments/assets/0cae7e3c-210a-4d5a-a189-67aa90daa982)



