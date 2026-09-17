# Domain Context

## Objetivo
Conter regras de negócio puras e invariantes financeiras.

## Entidades e Value Objects
- `Money` (imutável, sem float, unidade mínima)
- `Wallet` (aggregate root, sem saldo negativo)
- `WagerTransaction` (tipos/estados com transições válidas)

## Restrições da camada
- Sem dependência de framework, banco, HTTP ou SDKs externos.
- Foco em validações, comportamento e segurança monetária.

## Features relacionadas
- `internal/domain/errors/errors.go`
- `internal/domain/money/money.go`
- `internal/domain/wallet/wallet.go`
- `internal/domain/transaction/transaction.go`
- `internal/domain/outbox/outbox.go`

## Processo otimizado para agente de IA

1. Comece pelo núcleo do domínio e valide invariantes antes de qualquer detalhe de infraestrutura. O valor monetário, a carteira e a transação são a base do sistema.
2. Modelar `Money` como imutável e livre de float. A validação e a aritmética devem impedir valores inválidos, overflow e moedas incompatíveis.
3. Faça da `Wallet` a fonte da regra de negócio para débito/crédito, com saldo nunca negativo e controle de concorrência explícito.
4. Estruture a máquina de estados de `WagerTransaction` com transições bem definidas e estados terminalizados para rejeições e falhas permanentes.
5. Preserve a imutabilidade do ledger e a semântica de “append-only”; correções devem ser representadas por novas entradas compensatórias, não por edição.
6. Sempre tratar erros de negócio como tipos de domínio e nunca como `panic`. Use contexto para cancelamento e operações I/O controladas.
7. Antes de escrever infra, confirme que o domínio prova todos os casos de borda: overflow, moeda incompatível, saldo insuficiente, transação duplicada e referência pendente.
