#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

docker compose \
  --env-file "$PROJECT_ROOT/.env" \
  -f "$PROJECT_ROOT/tools/postgresql/compose.yaml" \
  down
