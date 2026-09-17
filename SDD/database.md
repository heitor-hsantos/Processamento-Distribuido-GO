# Database Context

## Objetivo
Centralizar requisitos de banco, migrações e integridade para persistência financeira.

## Contexto
- Banco principal: PostgreSQL.
- Acesso: `github.com/jackc/pgx/v5`.
- Migração inicial: `migrations/001_init_schema.sql`.

## Regras de Integridade
- Proibir saldo negativo com constraints (ex.: `CHECK (balance_cents >= 0)`).
- Garantir idempotência com persistência da chave e deduplicação por hash de payload.
- Preservar consistência para operações concorrentes.

## Features relacionadas
- `migrations/001_init_schema.sql`
- `internal/infrastructure/database/repository/wallet_repository.go`
- `internal/infrastructure/database/repository/transaction_repository.go`
- `internal/infrastructure/database/transaction/unit_of_work.go`

## Observações
- O banco deve reforçar invariantes do domínio, não apenas a aplicação.

## Processo otimizado para agente de IA

1. Modelar invariantes antes de implementar o código de acesso. Saldo negativo, duplicidade e ledger imutável devem existir no banco e no domínio.
2. Defina a estrutura de tabelas e constraints com foco em consistência financeira: `CHECK (balance_cents >= 0)`, unicidade de transações, chaves de idempotência e transações associadas ao wallet.
3. Especifique claramente as fronteiras de transação: saldo, ledger e outbox devem ser persistidos no mesmo `BEGIN/COMMIT` para evitar estados incongruentes.
4. Use `pgx`/SQL explícito e mantenha operações atômicas, évitando lógica de saldo em memória que possa divergir do banco em execução concorrente.
5. Sempre planeje a recuperação após falhas: outbox, inbox e retries devem registrar o estado final no banco para reprocessamento seguro.
6. Gere testes de concorrência e de idempotência no nível do banco para validar que duas operações conflitantes não geram saldo inconsistente.
7. Documente a estratégia de lock e a semântica de versionamento para explicar por que cada operação é segura em cenário distribuído.
