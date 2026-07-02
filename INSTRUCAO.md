# Instruções de Execução

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)

## Clone do Repositório

```bash
git clone https://github.com/lucasterencio/teste-golang-api.git
cd teste-golang-api
```

## Configuração

### Variáveis de Ambiente

O projeto utiliza o arquivo `.env` para configuração. As variáveis necessárias são:

```env
# Banco de Dados
DB_HOST=dbhost
DB_PORT=dbport
DB_USER=dbuser
DB_PASSWORD=dbpass
DB_NAME=dbname

# Aplicação
PORT=port

# PostgreSQL (para o container do banco)
POSTGRES_USER=dbuser
POSTGRES_PASSWORD=dbpass
POSTGRES_DB=dbname
```

## Executando com Docker (Recomendado)

### 1. Subir a aplicação completa

```bash
docker compose up --build
```

A aplicação estará disponível em: `http://localhost:3000`

### 2. Rodar em background

```bash
docker compose up -d --build
```

### 3. Parar a aplicação

```bash
docker compose down
```


### Rodar todos os testes

```bash
go test ./... -v
```

### Rodar testes de um pacote específico

```bash
# Testes dos handlers
go test ./src/handlers/... -v

# Testes com cobertura
go test ./... -cover
```

## Endpoints da API

### Swagger

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | `/swagger/index.html` | Acessar e testar todas as rotas disponíveis |
