# Mensageria Context

## Objetivo
Estruturar consumo/publicação de eventos com confiabilidade para processamento financeiro.

## Contexto
- Infra prevista: SQS (AWS SDK v2) / LocalStack no ambiente local.
- Workers em `internal/infrastructure/messaging`.

## Regras
- Inbox para deduplicação de mensagens recebidas.
- Outbox para garantir publicação após commit.
- Evitar perda de evento em falha de infraestrutura.

## Features relacionadas
- `internal/infrastructure/messaging/consumer/inbox_worker.go`
- `internal/infrastructure/messaging/consumer/bet_consumer.go`
- `internal/infrastructure/messaging/publisher/outbox_publisher.go`

## Variáveis de ambiente relacionadas
- `SQS_QUEUE_URL`

## Processo otimizado para agente de IA

1. A mensageria deve funcionar como mecanismo de entrega confiável, não como fonte de verdade financeira. O banco continua sendo a autoridade final.
2. Defina a arquitetura de filas com `inbox` para deduplicação e `outbox` para publicação pós-commit, preservando consistência em cenários de falha.
3. Sempre publique eventos somente após o commit da transação que alterou o domínio, nunca antes, para evitar inconsistência entre saldo e notificação.
4. O consumidor deve ser idempotente, tolerar retry e marcar mensagens como concluídas só quando o processamento financeiro estiver persistido com sucesso.
5. Implemente mecanismos de DLQ, retries com backoff e reprocessamento por worker após reinicialização.
6. Teste cenários de mensagem duplicada, ordem invertida (rollback antes da transação), falha de publicação e reconciliação de eventos pendentes.
7. Quando o worker falhar, mantenha o registro do status no banco para reprocessamento seguro e sem perda de eventos.
