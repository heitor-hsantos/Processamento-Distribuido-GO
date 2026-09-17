# Jungle Gaming

Este repositório é um scaffold inicial para um backend de apostas em Go seguindo DDD, Domain-Driven Design, com foco em carteira, ledger e processamento distribuído.

## Estrutura principal

- `cmd/api` — ponto de entrada da aplicação
- `internal/domain` — núcleo de negócio
- `internal/application` — casos de uso e portas
- `internal/infrastructure` — adaptadores externos e HTTP
- `internal/di` — configuração do Uber Fx
- `migrations` — scripts SQL para PostgreSQL

## Executar

```bash
make tidy
make build
make run
```

## Infra local

```bash
docker compose up -d
```

## Health checks e observabilidade

A aplicação expõe os endpoints de live/ready em:

```bash
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

Além disso, a camada de infraestrutura fornece logs estruturados em JSON e um contador simples de métricas via `internal/infrastructure/observability` para rastrear resultados operacionais, problemas de mensageria e sincronização financeira.

## Concorrência e recuperação

A lógica de carteira e transações foi organizada para permitir locks por linha em PostgreSQL e garantir que a versão do agregado seja incrementada em alterações fiscais. O fluxo de SQS/Outbox/InBox mantém deduplicação por `messageId` e hash do payload para evitar reprocessamentos duplicados.

## Próximos passos

- Implementar SQL real com pgx e transações ACID
- Adicionar repositórios concretos e workers SQS
- Finalizar autenticação OIDC e middleware de idempotência
- Adicionar testes de concorrência e fluxo financeiro
