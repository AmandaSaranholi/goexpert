### Google cloud Run: Weather API

Esse projeto implementa um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema foi publicado no Google Cloud Run.

### Requisitos

Go 1.21.6
Docker
Docker Compose


### Configuração do Ambiente

1. **Clone o repositório**

   ```sh
   git clone git@github.com:AmandaSaranholi/goexpert.git
   cd google-cloud-run
   ```

2. **Crie e configure o aquivo .env**

   ```sh
   cp .env.example .env
   # Modifique o arquivo .env como desejar
   ```

3. **Inicie o servidor Go com Docker Compose**

   ```sh
   docker-compose up -d
   ```

4. **O servidor vai responder no endereço `http://localhost:8080`**

5. **Faça requisições http para a URL do servidor usando umo navegador, HTTP Client ou até mesmo cURL**

   Exemplo acessando navegador:

   ```sh
    http://localhost:8080/weather/99999999    
   ```

6. **Execute os testes**

   ```sh
   docker compose exec app go test -v ./test
   ```


### URL do Projeto rodando na Google Cloud Run

Basta acessar a URL abaixo e substituir o zipcodepelo CEP desejado:
https://cloudrun-goexpert-rgioguuwia-uc.a.run.app/weather/{zipcode}


### Exemplos

**CEP Existente**
https://cloudrun-goexpert-rgioguuwia-uc.a.run.app/weather/17280276


**CEP Inexistente**
https://cloudrun-goexpert-rgioguuwia-uc.a.run.app/weather/17280275


**CEP Inválido**
https://cloudrun-goexpert-rgioguuwia-uc.a.run.app/weather/123
