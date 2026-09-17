# Idempotência Context

## Objetivo
Garantir que reenvios não dupliquem efeitos financeiros.

## Regras obrigatórias
- Persistir chave de idempotência.
- Usar hash de payload para deduplicação.
- Requisição repetida deve reutilizar resposta quando aplicável, com sinalização de replay.

## Superfícies
- HTTP middleware de idempotência.
- Persistência transacional para deduplicação.
- Fluxos de mensageria via inbox.

## Features relacionadas
- `internal/infrastructure/http/middleware/idempotency.go`
- `internal/infrastructure/messaging/consumer/inbox_worker.go`
- `internal/infrastructure/database/repository/transaction_repository.go`

## Processo otimizado para agente de IA

1. Trate idempotência como requisito de negócio e de infraestrutura, não como otimização de rede. O sistema deve reconhecer reenvios persistentes e reutilizar a resposta correta.
2. Persistir a chave de idempotência e o hash do payload antes de qualquer efeito financeiro, de forma transacional e durável.
3. Para HTTP e SQS, a deduplicação deve distinguir replays válidos de entradas de transação real ou duplicada por falha de transporte.
4. O replay deve devolver a mesma resposta original quando a operação já foi processada, além de rótulos como `idempotentReplay` para observabilidade.
5. Combine idempotência com inbox/outbox para garantir que uma mensagem duplicada não crie um novo efeito financeiro mesmo quando recebida mais de uma vez.
6. Faça testes de paralelismo com 50+ mensagens iguais para provar que o sistema produz um único ledger e mantém o estado auditable.
7. Toda implementação de idempotência deve considerar reinicialização, crash após commit e retries em workers.
