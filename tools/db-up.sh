#!/bin/bash

set -e

# Project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

docker compose \
  --env-file "$PROJECT_ROOT/.env" \
  -f "$PROJECT_ROOT/tools/postgresql/compose.yaml" \
  up -d
