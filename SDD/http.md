# HTTP Context

## Objetivo
Expor APIs da aplicação com handlers e middlewares alinhados às regras de negócio.

## Contexto
- Handlers e roteamento em `internal/infrastructure/http`.
- Middlewares para autenticação e idempotência.

## Regras
- HTTP apenas adapta entrada/saída.
- Regras de negócio ficam no domínio/aplicação.

## Features relacionadas
- `internal/infrastructure/http/handler/wallet.go`
- `internal/infrastructure/http/handler/transaction.go`
- `internal/infrastructure/http/middleware/auth.go`
- `internal/infrastructure/http/middleware/idempotency.go`
- `internal/infrastructure/http/router/router.go`

## Processo otimizado para agente de IA

1. Mantenha a camada HTTP como adaptadora: receber DTOs, validar entrada básica e delegar ao caso de uso correto.
2. Use middlewares para autenticação, idempotência e rastreio de correlação; não misture lógica de negócio ou persistência nos handlers.
3. Defina contratos de request e response com payloads explícitos, incluindo `Idempotency-Key`, `providerId`, `currency` e identificadores externos.
4. Ao falhar, retorne erro consistente com status HTTP apropriado e detalhe de negócio sem expor detalhes internos do domínio.
5. Sempre que houver replays ou duplicidade, a resposta deve indicar de forma clara se o resultado foi reutilizado por idempotência.
6. Teste apis em cenários de concorrência, autorização e reprocessamento para garantir que a camada HTTP não introduza inconsistência.
7. Escreva rotas e handlers com uma intenção clara: `wallets`, `transactions` e reconciliation devem ser observáveis e sem lógica escondida em middleware.
