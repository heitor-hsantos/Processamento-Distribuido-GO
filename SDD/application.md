# Application Context

## Objetivo
Orquestrar casos de uso através de ports, sem acoplamento à infraestrutura.

## Contexto
- Define contratos em `internal/application/port`.
- Implementa casos de uso em `internal/application/usecase`.

## Regras
- Depender apenas do domínio e interfaces.
- Não conhecer detalhes de PostgreSQL, HTTP, SQS ou OIDC.

## Features relacionadas
- `internal/application/port/repository.go`
- `internal/application/usecase/create_wallet.go`
- `internal/application/usecase/process_bet.go`
- `internal/application/usecase/reconcile.go`

## Processo otimizado para agente de IA

1. Comece pela regra de negócio do caso de uso e só depois pela infraestrutura. A aplicação deve expressar intenção de negócio, não detalhes de SQL, HTTP, SQS ou OIDC.
2. Defina os ports antes da implementação concreta. Cada caso de uso deve depender de interfaces pequenas e explícitas, evitando acoplamento com PostgreSQL, AWS ou frameworks.
3. Mantenha a camada de aplicação como orquestradora: criar carteira, processar aposta, reconciliar saldo e publicar resultados via eventos sem esconder regras de domínio.
4. Sempre valide a sequência: domínio -> ports -> use cases -> infraestrutura. Isso reduz retrabalho e ajuda a manter testes e integrações coerentes.
5. Construa testes de comportamento antes de código de infraestrutura, cobrindo saldo insuficiente, transação rejeitada, rejeição de moeda incompatível e evolução de estado.
6. Quando o caso de uso precisar de transação, concentre o contrato de unidade de trabalho no port apropriado e deixe a persistência resolver a implementação final.
7. Ao gerar código, priorize nomes explícitos, retornos de erro tipados e tratamento de `context.Context` para permitir cancelamento, observabilidade e concorrência controlada.
