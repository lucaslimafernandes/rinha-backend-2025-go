# 🐔 Rinha de Backend - 2025 (Go)

Autor: Lucas Lima Fernandes

[Github](https://github.com/lucaslimafernandes/)
[Linkedin](https://www.linkedin.com/in/lucaslimafernandes/)

Este projeto é a minha submissão para a terceira edição da [Rinha de Backend 2025](https://github.com/zanfranceschi/rinha-de-backend-2025), um desafio que visa aprendizado e compartilhamento de conhecimento na comunidade backend.

## 🚀 Descrição

O desafio de 2025 consiste em intermediar pagamentos entre dois serviços de processadores com diferentes taxas e confiabilidades. O sistema precisa:

- Verificar a disponibilidade dos processadores periodicamente
- Roteá-los com base em regras inteligentes
- Registrar os pagamentos com persistência
- Disponibilizar endpoints administrativos e analíticos

Minha solução foi escrita em **Go** com uso de **PostgreSQL**, **Docker**, **NGINX** e arquitetura leve, priorizando estabilidade e clareza.

---

## 🧱 Estrutura do Projeto

```
.
├── main.go                    # Entrypoint da aplicação
├── models/                    # Lógica de banco de dados
│   ├── database.go
│   └── insert.go
├── utils/                     # Utilitários e lógica de negócio
│   ├── payments.go
│   ├── utils.go
│   └── worker.go
├── nginx.conf                 # Configuração de load balancer NGINX
├── docker-compose.yaml        # Orquestração dos serviços
├── Dockerfile                 # Build da aplicação Go
├── go.mod / go.sum            # Dependências Go
└── README.md
```

---

## ⚙️ Como Executar

### 1. Pré-requisitos

- Docker e Docker Compose instalados
- Executar serviços dos `payment processors`, [instruções de execução do rinha](https://github.com/zanfranceschi/rinha-de-backend-2025/blob/main/rinha-test/README.md)

### 2. Subir os serviços

```bash
docker-compose up --build
```

Isso irá:
- Criar dois containers Go (`go_rinha_1` e `go_rinha_2`)
- Criar um banco de dados PostgreSQL
- Subir o NGINX como load balancer

A aplicação estará disponível em `http://localhost:9999`.

---

## 🧪 Endpoints

### `POST /payments`

Intermedia um novo pagamento. Roteia automaticamente para o processador `default` ou `fallback` com base na disponibilidade.

```json
{
  "id": 1,
  "valor": 1500,
  "moeda": "BRL"
}
```

### `GET /payments-summary?from=2025-07-01T00:00:00Z&to=2025-07-31T23:59:59Z`

Retorna um resumo estatístico dos pagamentos realizados no intervalo.

### `POST /admin/purge-payments`

Apaga os registros de pagamentos salvos (modo admin).

### `GET /healthy`

Retorna o status de saúde da aplicação.

---

## 🛠️ Tecnologias Utilizadas

- [Go 1.21+](https://golang.org/)
- [Docker + Docker Compose](https://docs.docker.com/)
- [PostgreSQL 14](https://www.postgresql.org/)
- [NGINX Alpine](https://hub.docker.com/_/nginx)

---

## 📈 Estratégia de Roteamento de Pagamentos

A lógica considera três cenários:

1. **default saudável** → envia para `default`
2. **default fora** → fallback saudável → envia para `fallback`
3. **ambos indisponíveis** → tenta `default` por último recurso

A verificação de saúde é feita a cada 5 segundos.

---

## 🧼 Limites de Recursos

Todos os serviços estão com limites de CPU e memória pré-definidos no `docker-compose.yaml`, conforme regras da Rinha.

| Serviço          | CPU (cores) | Memória (MB) |
|------------------|-------------|---------------|
| db_pg            | 0.5         | 200           |
| nginx            | 0.3         | 50            |
| go_rinha_1       | 0.35        | 50            |
| go_rinha_2       | 0.35        | 50            |
| **Total**        | **1.5**     | **350**       |



---

## 📣 Créditos

Desenvolvido por [Lucas Lima Fernandes](https://github.com/lucaslimafernandes)  
Para a comunidade da [Rinha de Backend](https://github.com/zanfranceschi/rinha-de-backend-2025) 

---

## 🐔 Links úteis

- [Prévia dos resultados - Kauê Fraga](https://rinhers.kauefraga.dev/)
- [Ranking em tempo real - Anderson Gomes](https://rinha2025.andersongomes.dev.br/)
- [Repositório oficial](https://github.com/zanfranceschi/rinha-de-backend-2025)
