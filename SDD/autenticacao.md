# Autenticação Context

## Objetivo
Preparar autenticação OIDC para proteger endpoints e identidade de chamadas.

## Contexto
- Estrutura preparada para OIDC.
- Integração local prevista com Keycloak.

## Features relacionadas
- `internal/infrastructure/auth/oidc.go`
- `internal/infrastructure/http/middleware/auth.go`

## Variáveis de ambiente relacionadas
- `OIDC_ISSUER_URL`
- `OIDC_AUDIENCE`

## Processo otimizado para agente de IA

1. Entenda que a autenticação é a fronteira externa da aplicação: o código deve validar o token antes de qualquer negócio financeiro ser executado.
2. Manter o providerId do contexto igual ao providerId do token; nunca aceitar uma requisição sem identificação explícita do provedora/cliente autorizado.
3. Defina os fluxos de autenticação em middleware HTTP e em clientes de infraestrutura, sem misturar validação OIDC com regras de saldo ou ledger.
4. Sempre separar autenticação, autorização e validação de domínio. O primeiro responde “quem é?”, o segundo “tem permissão?”, e o terceiro “é válido financeiramente?”.
5. Ao implementar, considere `client_credentials` para serviços, validação de issuer/audience, `exp`, `iat` e rejeição de tokens inválidos, expirados ou sem providerId.
6. Escreva testes para ausência de token, token inválido, providerId divergente e fluxo de request autenticado válido.
7. Quando a autenticação falhar, a resposta deve ser clara e determinística, sem que o processamento de negócios avance em um estado ambíguo.
