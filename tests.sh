#!/usr/bin/env bash
# Encerra o script caso algum comando falhe e garante que erros em pipes sejam pegos
set -e
set -o pipefail

# Cores para o output do terminal
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Iniciando pipeline de testes do Processamento Distribuído (Go) ===${NC}\n"

# 0. Instalação do formatador visual (gotestsum)
echo -e "${GREEN}[0/5] Verificando dependências visuais (gotestsum)...${NC}"
if ! command -v gotestsum &> /dev/null; then
  echo -e "${YELLOW}gotestsum não encontrado. Instalando temporariamente para uma melhor visualização...${NC}"
  go install gotest.tools/gotestsum@latest
  # Garante que o binário do Go esteja no PATH para a execução atual
  gopath=$(go env GOPATH)
  export PATH="$PATH:$gopath/bin"
fi
echo -e "✔ gotestsum pronto para uso.\n"

# 1. Verificação Estática e Formatação
echo -e "${GREEN}[1/5] Executando formatação e linting (gofmt e go vet)...${NC}"
gofmt -s -w .
go vet ./...
echo -e "✔ Linting concluído com sucesso.\n"

# 2. Testes Unitários Isolados (Domínio)
echo -e "${GREEN}[2/5] Executando testes unitários (Domínio e Invariantes)...${NC}"
gotestsum --format testname -- -race ./internal/domain/...
echo -e "✔ Testes unitários do domínio concluídos.\n"

# 3. Inicialização da Infraestrutura de Integração
echo -e "${GREEN}[3/5] Subindo dependências de infraestrutura (PostgreSQL, SQS, Keycloak)...${NC}"
docker compose up -d --build
echo -e "${YELLOW}Aguardando o PostgreSQL processar as migrations e ficar pronto...${NC}"
# Tenta conectar no banco a cada 2 segundos, até um limite de 30 tentativas
# Ajuste "postgres", "postgres" e "jungle_gaming" conforme seu docker-compose.yml atual
for i in {1..30}; do
  if docker compose exec -T postgres pg_isready -U postgres -d jungle_gaming > /dev/null 2>&1; then
	echo -e "✔ PostgreSQL está pronto e migrations foram aplicadas!"
	break
  fi
  echo -e "Aguardando banco de dados... ($i/30)"
  sleep 2
  if [ "$i" -eq 30 ]; then
	echo -e "${RED}Erro: Tempo limite atingido aguardando serviços.${NC}"
	docker compose logs
	docker compose down -v
	exit 1
  fi
done

# 4. Testes de Integração e Concorrência (Aplicação e Infraestrutura)
echo -e "${GREEN}[4/5] Executando testes de integração, concorrência e resiliência...${NC}"
echo -e "${YELLOW}-> Execução Padrão (inclui testes de integração quando INTEGRATION_TEST=1):${NC}"
# Só tenta rodar testes de integração a partir do host se a porta 5432 estiver exposta no host
if command -v nc &>/dev/null && nc -z localhost 5432; then
  echo -e "${GREEN}Postgres acessível em localhost:5432 — executando testes com integração habilitada.${NC}"
  INTEGRATION_TEST=1 gotestsum --format testname -- ./...
  echo -e "\n${YELLOW}-> Execução com Race Detector (Obrigatório para validar Concorrência):${NC}"
  INTEGRATION_TEST=1 gotestsum --format testname -- -race ./...
  echo -e "✔ Testes de integração e concorrência concluídos.\n"
else
  echo -e "${YELLOW}Postgres não acessível em localhost:5432 — pulando testes de integração que requerem conexão do host.${NC}"
  echo -e "${YELLOW}Rodando suíte normal (sem integração host->DB).${NC}"
  gotestsum --format testname -- ./...
  echo -e "\n${YELLOW}-> Execução com Race Detector:${NC}"
  gotestsum --format testname -- -race ./...
  echo -e "✔ Testes de integração e concorrência concluídos.\n"
fi
docker compose down -v
echo -e "✔ Ambiente limpo.\n"
echo -e "${GREEN}======================================================${NC}"
echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM! A solução está aderente.${NC}"
echo -e "${GREEN}======================================================${NC}"