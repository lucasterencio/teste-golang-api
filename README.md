# 🚀 Teste Técnico - API de Reserva de Salas

## 🎯 Objetivo

Desenvolver uma API REST para gerenciamento de reservas de salas.

> **Observação:** A modelagem das entidades, arquitetura da aplicação e decisões de implementação ficam livres para o candidato.

---

## 🛠️ Requisitos Obrigatórios

A solução deve utilizar:

* ✅ Go
* ✅ Gin
* ✅ PostgreSQL
* ✅ Docker
* ✅ Docker Compose
* ✅ Testes unitários
* ✅ Paginação (`page` e `limit`)
* ✅ Filtros básicos
* ✅ Tratamento de erros personalizados
* ✅ README com instruções para execução

---

## ⭐ Diferenciais (Bônus)

Os itens abaixo não são obrigatórios, mas serão considerados como diferencial:

* 📖 Swagger / OpenAPI
* 🗄️ Migrations
* 🧪 Cobertura mínima de testes de 70%
* 🔄 CI com GitHub Actions
* 📋 Logs estruturados
* ⚙️ Configuração via arquivo `.env`
* 🚀 Cache implementado com Redis
* 📊 Monitoramento básico de métricas e health check

---

## 📌 Exemplos de Requisições

### 📄 Paginação

```http
GET /reservations?page=1&limit=10
```

### 🔎 Filtros

Filtrar reservas por sala:

```http
GET /reservations?room_id=1
```

Filtrar reservas por responsável:

```http
GET /reservations?responsible=joao
```

---

## 💡 O que será avaliado

* Organização e estrutura do projeto
* Qualidade do código
* Boas práticas em Go
* Modelagem de dados
* Tratamento de erros
* Cobertura e qualidade dos testes
* Uso adequado do PostgreSQL
* Utilização de Docker e Docker Compose
* Estratégia de cache (caso implementado)
* Clareza da documentação

Boa sorte! 🚀
