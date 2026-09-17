# Migrations Feature

## Objetivo
Versionar estrutura de banco e regras de integridade desde o início do scaffold.

## Contexto
- Migrações ficam em `migrations/`.
- O estado atual menciona a migração `001_init_schema.sql`.

## Requisitos
- Criar estruturas iniciais para carteiras, transações e suporte a consistência.
- Incluir constraints financeiras obrigatórias.
- Preparar base para inbox/outbox e evolução incremental.

## Feature relacionada
- `migrations/001_init_schema.sql`

## Processo otimizado para agente de IA

1. Comece pela migração como contrato de integridade do sistema: schema, constraints, índices e unicidade devem refletir as regras financeiras.
2. Defina a estrutura de banco em etapas: carteira, transação, ledger, inbox e outbox. Cada tabela deve representar uma invariável do domínio.
3. Especifique constraints que impeçam saldo negativo, transação duplicada e ledger editável, validando as regras antes de mover para código de aplicação.
4. As migrações devem permitir evolução incremental com versionamento e documentação de rollback, favorecendo manutenção em ambiente distribuído.
5. Sempre que um caso de uso requer consistência transacional, confirme que a tabela e as regras do banco suportam a mesma garantia percebida pela aplicação.
6. Teste as migrações em execução real, validando criação, deduplicação, checks e execução concorrente para não gerar estado inválido em produção.
7. Documente cada migration com o motivo da mudança, deixando claro quais regras financeiras e de auditoria ela estabiliza.
