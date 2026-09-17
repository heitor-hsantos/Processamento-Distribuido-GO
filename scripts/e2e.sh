#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

AWS_ENDPOINT="${AWS_ENDPOINT:-http://localhost:4566}"
AWS_REGION="${AWS_REGION:-us-east-1}"
QUEUE_NAME="${QUEUE_NAME:-wager-transactions.fifo}"
DLQ_NAME="${DLQ_NAME:-wager-transactions-dlq.fifo}"
DB_CONTAINER="${DB_CONTAINER:-jungegaming-test-postgres-1}"
POSTGRES_USER="${POSTGRES_USER:-postgres}"
POSTGRES_DB="${POSTGRES_DB:-jungle_gaming}"

echo "[1/8] Subindo infraestrutura..."
docker compose up --build -d

echo "[2/8] Aplicando migrations..."
for migration in migrations/*.sql; do
  echo "  -> $(basename "$migration")"
  docker compose exec -T postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "/workspace/$migration" 2>/dev/null || \
  docker compose exec -T postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "/$migration"
done

echo "[3/8] Criando DLQ FIFO..."
aws --endpoint-url="$AWS_ENDPOINT" --region "$AWS_REGION" sqs create-queue \
  --queue-name "$DLQ_NAME" \
  --attributes FifoQueue=true,ContentBasedDeduplication=false >/dev/null

echo "[4/8] Lendo ARN da DLQ..."
DLQ_URL="$(aws --endpoint-url="$AWS_ENDPOINT" --region "$AWS_REGION" sqs get-queue-url --queue-name "$DLQ_NAME" --query QueueUrl --output text)"
DLQ_ARN="$(aws --endpoint-url="$AWS_ENDPOINT" --region "$AWS_REGION" sqs get-queue-attributes --queue-url "$DLQ_URL" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)"

echo "[5/8] Criando fila principal FIFO com redrive..."
aws --endpoint-url="$AWS_ENDPOINT" --region "$AWS_REGION" sqs create-queue \
  --queue-name "$QUEUE_NAME" \
  --attributes "FifoQueue=true,ContentBasedDeduplication=false,RedrivePolicy={\"deadLetterTargetArn\":\"$DLQ_ARN\",\"maxReceiveCount\":\"5\"}" >/dev/null

echo "[6/8] Executando testes padrão..."
go test ./...

echo "[7/8] Executando testes com race detector..."
go test -race ./internal/domain/... ./internal/application/...

echo "[8/8] Executando go vet..."
go vet ./...

echo "E2E bootstrap concluído."
echo "Filas criadas:"
echo "  - $QUEUE_NAME"
echo "  - $DLQ_NAME"
echo "Próximo passo manual: provisionar realm/client no Keycloak e executar cenários autenticados."
