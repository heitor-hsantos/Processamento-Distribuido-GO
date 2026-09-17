# Repository Context

## Objetivo
Definir responsabilidades e contratos dos repositórios na arquitetura em camadas (DDD + Ports and Adapters).

## Contexto Arquitetural
- A camada `internal/application` define interfaces em `internal/application/port`.
- A camada `internal/infrastructure/database/repository` implementa essas interfaces.
- A camada de aplicação não conhece PostgreSQL/pgx diretamente.

## Regras
- Repositórios devem ser acessados por interfaces (ports).
- Implementações concretas vivem em infraestrutura.
- Conversões de formatos externos e persistência ficam fora do domínio.
- Consistência deve respeitar regras financeiras e transacionais.

## Features relacionadas
- `internal/application/port/repository.go`
- `internal/infrastructure/database/repository/wallet_repository.go`
- `internal/infrastructure/database/repository/transaction_repository.go`
- `internal/infrastructure/database/transaction/unit_of_work.go`

## Observações
- O `Unit of Work` coordena fronteira transacional para casos de uso financeiros.

## Processo otimizado para agente de IA

1. A camada de repositório deve implementar interfaces de aplicação, não conter regras de negócio nem lógica de HTTP ou mensageria.
2. Cada operação deve mapear claramente entrada de domínio para SQL e vice-versa, com foco em atomicidade e controle de concorrência.
3. Organize os repositórios por agregado: carteira, transação, ledger, inbox/outbox; use consultas SQL explícitas e verificáveis.
4. O `Unit of Work` deve agrupar leitura e escrita dentro de uma transação financeira, garantindo que saldo, ledger e outbox saiam do mesmo commit.
5. Implemente tratamento explícito para conflitos de atualização, retries e idempotência persistente para garantir reprocessamento seguro.
6. Produza testes que validem pesquisa por carteira, alteração de saldo, transações `BET/WIN/REFUND/ROLLBACK` e recuperação após falha parcial.
7. Ao alterar consultas, mantenha os contratos de interface estáveis e reduza a lógica de domínio embalado em SQL para manter rastreabilidade e manutenção.
