# Processamento Distribuído de Apostas em Go

Este projeto implementa um serviço em **Go** utilizando **Uber Fx** para processar operações financeiras de provedores de jogos em um ambiente distribuído. A solução oferece uma API HTTP e um consumidor de mensagens SQS com garantias estritas de idempotência persistente, integridade do ledger (append-only) e controle de concorrência.

## 1. Decisões Técnicas e Arquiteturais

O projeto foi estruturado com base nos princípios de **Domain-Driven Design (DDD)** e **Clean Architecture**, garantindo que o domínio financeiro fique isolado de frameworks e da infraestrutura de I/O.

* **Composição e Ciclo de Vida:** O **Uber Fx** gerencia a inicialização e o *graceful shutdown* dos handlers HTTP e workers do SQS, assegurando que o processamento em andamento seja concluído antes da liberação de recursos.


* **Precisão Monetária:** O tipo `Money` foi modelado como um Value Object utilizando `int64` em unidades mínimas, evitando flutuações e erros de arredondamento inerentes a tipos `float`.


* **Controle de Concorrência:** Implementado via **Unit of Work**, o controle de concorrência ocorre na raiz do agregado (carteira) com *row-level locking* no PostgreSQL. Isso permite a atualização atômica do saldo e a inserção no ledger sem gerar gargalos globais.


* **Idempotência e Mensageria (Inbox/Outbox):**
* Rejeição de duplo processamento garantida por idempotência persistente (middleware HTTP) e pelo padrão **Inbox** (SQS).


* Eventos são publicados via **Transactional Outbox**, sendo gravados na mesma transação SQL da aposta e despachados assincronamente por um worker dedicado com suporte a *retry*.




* **Autenticação Multi-Tenant:** Integração com IdP OIDC (OAuth 2.0). O `providerId` derivado do token restringe rigorosamente o acesso apenas às transações do próprio provedor.


* **Documentação em SDD:** Regras de domínio, banco de dados, mensageria e autenticação estão documentadas na pasta `/SDD`, otimizando o contexto para manutenção contínua e integração com agentes de IA.



---

## 2. Pré-requisitos

Para executar e testar o projeto a partir de um *checkout* limpo, você precisará de:

* **Go** (versão mais recente declarada no `go.mod`)


* **Docker** e **Docker Compose**


## 3. Configuração e Inicialização

```

Suba a infraestrutura completa (PostgreSQL, LocalStack SQS e Keycloak) via Docker Compose:

```bash
docker compose up --build

```

* **Migrations:** Aplicadas automaticamente na primeira inicialização do PostgreSQL a partir do diretório `/migrations`.


* **Filas SQS:** As filas `wager-transactions.fifo` e `wager-transactions-dlq.fifo` são criadas via script de inicialização do LocalStack.


* **Keycloak:** O ambiente é provisionado com a identidade `provider-a` para testes via `client_credentials`.



## 4. Executando a Aplicação

Com a infraestrutura ativa, inicie a aplicação:

```bash
go run cmd/api/main.go

```

## 5. Testes e Validação (E2E)

A resiliência contra dupla cobrança e falhas de concorrência é validada pela suíte de testes. Utilize o script de execução automatizada para validar o projeto completo com o **Race Detector** ativo:

```bash
chmod +x tests.sh
./tests.sh

```

*(Para testes manuais locais, utilize `go test -v -race ./...`)*

## 6. Exemplos de Chamadas da API

### Autenticação (Obtendo o Token)

```http
POST http://localhost:8080/realms/jungle-gaming/protocol/openid-connect/token
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&client_id=provider-a&client_secret=SECRETO_DO_ENV

```

### Criação de Carteira (Interno/Admin)

```http
POST /wallets
Content-Type: application/json
Authorization: Bearer <SEU_TOKEN>

{
  "playerId": "0192f28f-5dc0-7d58-bdb2-814ad6a0f4a1",
  "initialBalance": { "amount": "1000.00", "currency": "BRL" }
}

```

### Processamento de Aposta

```http
POST /wagering/transactions
Content-Type: application/json
Idempotency-Key: provider-a:transaction-123
Authorization: Bearer <SEU_TOKEN>

{
  "providerId": "provider-a",
  "externalTransactionId": "transaction-123",
  "playerId": "0192f28f-5dc0-7d58-bdb2-814ad6a0f4a1",
  "walletId": "0192f291-27dd-7d3f-8071-5f8685deef37",
  "roundId": "round-987",
  "gameId": "fortune-chimp",
  "kind": "BET",
  "money": { "amount": "25.00", "currency": "BRL" }
}

```

### Reconciliação

```http
POST /wallets/0192f291-27dd-7d3f-8071-5f8685deef37/reconciliation
Authorization: Bearer <SEU_TOKEN>

```
