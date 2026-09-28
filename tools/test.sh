#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Test Runner
# ============================================================

set -u

total_start_ms=$(date +%s%3N)

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

backend_start_ms=$(date +%s%3N)

backend_output="$(
    cd "${APP_DIR}" &&
    go test -v ./backend/... 2>&1
)"
backend_status=$?

backend_total=0
backend_ng=0
backend_classification_error=0

while IFS= read -r line; do
    if [[ "${line}" =~ ^---\ (PASS|FAIL):\ Test([BC])_(.+)\ \((.+)\)$ ]]; then
        result="${BASH_REMATCH[1]}"
        class="${BASH_REMATCH[2]}"
        name="${BASH_REMATCH[3]//_/ }"
        duration="${BASH_REMATCH[4]}"

        ((backend_total+=1))

        if [[ "${result}" == "PASS" ]]; then
            echo "[${class}] PASS   ${name} (${duration})"
        else
            echo "[${class}] FAIL   ${name} (${duration})"
            ((backend_ng+=1))
        fi

    elif [[ "${line}" =~ ^---\ (PASS|FAIL):\ (Test[^[:space:]]+)\ \((.+)\)$ ]]; then
        name="${BASH_REMATCH[2]}"
        duration="${BASH_REMATCH[3]}"

        ((backend_total+=1))
        ((backend_ng+=1))
        backend_classification_error=1

        echo "[?] ERROR  ${name} (${duration})"
    fi
done <<< "${backend_output}"

backend_end_ms=$(date +%s%3N)
backend_elapsed_ms=$((backend_end_ms - backend_start_ms))

echo


# ------------------------------------------------------------
# Frontend
# ------------------------------------------------------------

echo "=== Frontend Test ==="
echo

frontend_start_ms=$(date +%s%3N)

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

frontend_total=0
frontend_ng=0
frontend_classification_error=0

while IFS= read -r line; do
    if [[ "${line}" =~ \[([BC])\][[:space:]]+(.+)[[:space:]]+([0-9]+ms)$ ]]; then
        class="${BASH_REMATCH[1]}"
        name="${BASH_REMATCH[2]}"
        duration="${BASH_REMATCH[3]}"

        ((frontend_total+=1))

        echo "[${class}] PASS   ${name} (${duration})"

    elif [[ "${line}" =~ \>[[:space:]]+([^[][^[:space:]].*)[[:space:]]+([0-9]+ms)$ ]]; then
        name="${BASH_REMATCH[1]}"
        duration="${BASH_REMATCH[2]}"

        ((frontend_total+=1))
        ((frontend_ng+=1))
        frontend_classification_error=1

        echo "[?] ERROR  ${name} (${duration})"
    fi
done <<< "${frontend_output_clean}"

frontend_end_ms=$(date +%s%3N)
frontend_elapsed_ms=$((frontend_end_ms - frontend_start_ms))

echo

total_end_ms=$(date +%s%3N)
total_elapsed_ms=$((total_end_ms - total_start_ms))

# ------------------------------------------------------------
# Result
# ------------------------------------------------------------

format_seconds() {
    local ms=$1
    printf "%d.%01ds" $((ms / 1000)) $(((ms % 1000) / 100))
}

ng_rate() {
    local total=$1
    local ng=$2

    if (( total == 0 )); then
        printf "0.0%%"
    else
        awk -v ng="${ng}" -v total="${total}" \
            'BEGIN { printf "%.1f%%", ng * 100 / total }'
    fi
}

total_tests=$((backend_total + frontend_total))
total_ng=$((backend_ng + frontend_ng))

if (( backend_ng == 0 && backend_status == 0 )); then
    backend_result="PASS"
else
    backend_result="FAIL"
fi

if (( frontend_ng == 0 && frontend_status == 0 )); then
    frontend_result="PASS"
else
    frontend_result="FAIL"
fi

if (( total_ng == 0 &&
      backend_status == 0 &&
      frontend_status == 0 )); then
    total_result="PASS"
else
    total_result="FAIL"
fi

echo "=== Result ==="
echo

printf "Backend  : %-4s   %d tests / %d NG (%s)   %s\n" \
    "${backend_result}" \
    "${backend_total}" \
    "${backend_ng}" \
    "$(ng_rate "${backend_total}" "${backend_ng}")" \
    "$(format_seconds "${backend_elapsed_ms}")"

printf "Frontend : %-4s   %d tests / %d NG (%s)   %s\n" \
    "${frontend_result}" \
    "${frontend_total}" \
    "${frontend_ng}" \
    "$(ng_rate "${frontend_total}" "${frontend_ng}")" \
    "$(format_seconds "${frontend_elapsed_ms}")"

echo "------------------------------------------------"

printf "Total    : %-4s   %d tests / %d NG (%s)   %s\n" \
    "${total_result}" \
    "${total_tests}" \
    "${total_ng}" \
    "$(ng_rate "${total_tests}" "${total_ng}")" \
    "$(format_seconds "${total_elapsed_ms}")"


# ending
if [[ ${backend_status} -eq 0 &&
      ${frontend_status} -eq 0 &&
      ${backend_classification_error} -eq 0 &&
      ${frontend_classification_error} -eq 0 ]]; then
    exit 0
else
    exit 1
fi
