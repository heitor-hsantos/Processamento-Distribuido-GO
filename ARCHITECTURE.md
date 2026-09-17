# Arquitetura do Sistema

## Visão geral

O sistema é estruturado em camadas para separar domínio, aplicação e infraestrutura:

- `internal/domain` contém regras de negócio imutáveis e robustas, sem dependência de Fx, HTTP ou banco de dados.
- `internal/application` define casos de uso e interfaces de porta para repositórios e serviços externos.
- `internal/infrastructure` efetua a adaptação ao PostgreSQL, SQS, OAuth/OIDC e HTTP.
- `internal/di` centraliza a wiring com Uber Fx e lifecycle hooks.

## Estratégia financeira

- Toda movimentação financeira usa inteiros em centavos, nunca float.
- Saldo negativo é proibido por constraint no banco e validação no agregado Wallet.
- O Ledger é append-only e registros imutáveis; correções são feitas por transações compensatórias.

## Concorrência e consistência

- A aplicação separa transações de domínio e alterações em contas por unidade de trabalho.
- O uso de versões do agregado e testes de concorrência permite proteger o cenário de lost update.
- O outbox e inbox garantem consistência entre commit SQL e publicação de eventos em filas externas.

## Segurança

- A autenticação é tratada por OIDC e middleware de HTTP para validar o token e extrair o `providerId` do contexto.
- A idempotência é exigida por chave e hash do payload para evitar duplicação persistente de transações.

## Observações

Este scaffold estabelece a base arquitetural e serviços de infraestrutura necessários para evoluir para a implementação completa do desafio de apostas distribuídas.
