#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Test Runner
# ============================================================

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
APP_DIR="${PROJECT_ROOT}/app"
FRONTEND_DIR="${APP_DIR}/frontend"

backend_status=0
frontend_status=0


# ------------------------------------------------------------
# Backend
# ------------------------------------------------------------

echo "=== Backend Test ==="
echo

backend_output="$(
    cd "${APP_DIR}" &&
    go test -v ./backend/... 2>&1
)"
backend_status=$?

echo "${backend_output}" |
while IFS= read -r line; do
    if [[ "${line}" =~ ^---\ PASS:\ TestB_(.+)\ \((.+)\)$ ]]; then
        name="${BASH_REMATCH[1]//_/ }"
        duration="${BASH_REMATCH[2]}"
        echo "[B] PASS  ${name} (${duration})"

    elif [[ "${line}" =~ ^---\ PASS:\ TestC_(.+)\ \((.+)\)$ ]]; then
        name="${BASH_REMATCH[1]//_/ }"
        duration="${BASH_REMATCH[2]}"
        echo "[C] PASS  ${name} (${duration})"
    fi
done

echo


# ------------------------------------------------------------
# Frontend
# ------------------------------------------------------------

echo "=== Frontend Test ==="
echo

frontend_output="$(
    cd "${FRONTEND_DIR}" &&
    npm run test -- --reporter=verbose 2>&1
)"
frontend_status=$?

# ANSI escape sequenceを除去
frontend_output_clean="$(
    printf '%s\n' "${frontend_output}" |
    sed $'s/\033\\[[0-9;]*m//g'
)"

printf '%s\n' "${frontend_output_clean}" |
while IFS= read -r line; do
    if [[ "${line}" =~ \[([BC])\][[:space:]]+(.+)[[:space:]]+([0-9]+ms)$ ]]; then
        class="${BASH_REMATCH[1]}"
        name="${BASH_REMATCH[2]}"
        duration="${BASH_REMATCH[3]}"

        echo "[${class}] PASS  ${name} (${duration})"
    fi
done

echo


# ------------------------------------------------------------
# Result
# ------------------------------------------------------------

if [[ ${backend_status} -eq 0 && ${frontend_status} -eq 0 ]]; then
    exit 0
else
    exit 1
fi
