#!/usr/bin/env bash

set -e

# Project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

# Load local environment settings if present
if [ -f ".env" ]; then
  set -a
  source ".env"
  set +a
fi

# Redocly telemetry is always disabled.
export REDOCLY_TELEMETRY=off

# Show Redocly update notices only in ARIADNE development mode.
if [ "${ARIADNE_DEVELOPMENT:-false}" != "true" ]; then
  export REDOCLY_SUPPRESS_UPDATE_NOTICE=true
fi

OAS_ROOT="./dist/api/oas"

generate_redoc() {
  SERVICE="$1"
  OAS="${OAS_ROOT}/${SERVICE}/openapi.yaml"
  HTML="${OAS_ROOT}/${SERVICE}/redoc.html"
  EXTERNAL_DOCS="${OAS_ROOT}/${SERVICE}/docs"

  echo "=================================================="
  echo "Service: ${SERVICE}"
  echo "=================================================="

  if [ ! -f "${OAS}" ]; then
    echo "ERROR: OpenAPI file not found."
    echo "       ${OAS}"
    return 1
  fi

  echo "[1/3] Lint OpenAPI..."
if ! npx redocly lint "${OAS}"; then
  echo
  echo "ERROR: OpenAPI lint failed."

  if [ -f "${HTML}" ]; then
    rm "${HTML}"
    echo "Removed old ReDoc: ${HTML}"
  fi

  if [ -d "${EXTERNAL_DOCS}" ]; then
    rm -rf "${EXTERNAL_DOCS}"
    echo "Removed old externalDocs: ${EXTERNAL_DOCS}"
  fi

  return 1
fi

  echo
  echo "[2/3] Generate externalDocs..."
  node ./tools/generate-external-doc.mjs "${SERVICE}"

  echo
  echo "[3/3] Generate ReDoc..."
  npx redocly build-docs "${OAS}" --output="${HTML}"

  echo
  echo "OK: ${HTML}"
  echo
}

# 引数あり：指定サービスのみ
if [ $# -gt 0 ]; then
  generate_redoc "$1"
  exit 0
fi

# 引数なし：全サービス
echo "Generate ReDoc for all services..."
echo

FOUND=false

for DIR in "${OAS_ROOT}"/*/; do
  [ -d "${DIR}" ] || continue

  FOUND=true
  SERVICE=$(basename "${DIR}")
  generate_redoc "${SERVICE}"
done

if [ "${FOUND}" = false ]; then
  echo "ERROR: No services found under ${OAS_ROOT}"
  exit 1
fi

echo "All ReDoc files generated successfully."
