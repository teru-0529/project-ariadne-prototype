#!/usr/bin/env bash

set -e

OAS_ROOT="./dist/api/oas"

generate_redoc() {
  SERVICE="$1"
  OAS="${OAS_ROOT}/${SERVICE}/openapi.yaml"
  HTML="${OAS_ROOT}/${SERVICE}/redoc.html"

  echo "=================================================="
  echo "Service: ${SERVICE}"
  echo "=================================================="

  if [ ! -f "${OAS}" ]; then
    echo "ERROR: OpenAPI file not found."
    echo "       ${OAS}"
    return 1
  fi

  echo "[1/2] Lint OpenAPI..."
if ! npx redocly lint "${OAS}"; then
  echo
  echo "ERROR: OpenAPI lint failed."

  if [ -f "${HTML}" ]; then
    rm "${HTML}"
    echo "Removed old ReDoc: ${HTML}"
  fi

  return 1
fi

  echo
  echo "[2/2] Generate ReDoc..."
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
