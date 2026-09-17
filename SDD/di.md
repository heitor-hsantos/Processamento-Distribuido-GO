# DI Context

## Objetivo
Centralizar composição da aplicação com Uber Fx.

## Contexto
- `internal/di` contém módulos de wiring.
- `cmd/api/main.go` inicializa aplicação.

## Regras
- Apenas a camada DI conhece composição completa.
- Lifecycle hooks de HTTP e workers devem estar centralizados.

## Features relacionadas
- `internal/di/module.go`
- `internal/di/app_fx.go`
- `internal/di/infra_fx.go`
- `cmd/api/main.go`

## Processo otimizado para agente de IA

1. Trate a DI como a camada de composição da aplicação: monte dependências em módulos pequenos e explícitos, não em um único ponto de acoplamento.
2. Garanta que o domínio e a aplicação não dependam de `fx`, HTTP, SQL ou SQS; a composição deve apenas conectar implementações concretas às interfaces.
3. Use `fx.Module` para separar responsabilidades (infra, domínio, app, HTTP, mensageria), deixando a estrutura rastreável e fácil de testar.
4. Defina `fx.Lifecycle` para iniciar/encerrar workers, HTTP server e conexões do banco de forma ordenada e segura.
5. Ao criar dependências, prefira construtores simples e injetáveis, com zero lógica de negócio em módulos de wiring.
6. Teste a composição com build e execução local para confirmar que os módulos se conectam corretamente e que shutdown não deixa recursos abertos.
7. Em mudanças de arquitetura, sempre preservar a ordem de inicialização: configuração -> banco -> repositórios -> casos de uso -> handlers/workers -> lifecycle.
