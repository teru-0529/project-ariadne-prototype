#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Backend Test Runner
# ============================================================

set -u

GREEN='\033[92m'
RED='\033[91m'
RESET='\033[0m'


# ------------------------------------------------------------
# Initial
# ------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
APP_DIR="${PROJECT_ROOT}/app"

status=0

# INFO: パラメータ判定（Dテスト実施有無）
arg="${1:-}"
param=""
echo

if [[ "${arg}" == "--logic" ]]; then
    param="-logic"
else
    echo "--- D tests are skipped ---"
    echo
fi


# ------------------------------------------------------------
# Backend
# ------------------------------------------------------------

echo "=== Backend Test ==="
echo

start_ms=$(date +%s%3N)

# INFO: go test 実行
output="$(
    cd "${APP_DIR}" &&
    go test -v ./backend/... ${param} 2>&1
)"
status=$?

total=0
ng=0

while IFS= read -r line; do

    # INFO: Test開始行削除
    if [[ "${line}" =~ ^===\ RUN[[:space:]]+ ]]; then
        continue
    fi

    # INFO: SKIP行削除
    if [[ "${line}" =~ ^---\ SKIP:[[:space:]]+ ]]; then
        continue
    fi

    # INFO: test結果行
    if [[ "${line}" =~ ^---\ (FAIL:|PASS:)[[:space:]]+Test_([^[:space:]]+)[[:space:]]+(\(.+\)) ]]; then

        if [[ "${BASH_REMATCH[1]}" == "FAIL:" ]]; then
            printf "  ${RED}●${RESET} ${BASH_REMATCH[1]} ${BASH_REMATCH[2]//_/ } ${BASH_REMATCH[3]}\n"
            ((ng+=1))
        else
            printf "  ${GREEN}●${RESET} ${BASH_REMATCH[1]} ${BASH_REMATCH[2]//_/ } ${BASH_REMATCH[3]}\n"
        fi
        ((total+=1))
        continue
    fi

    # INFO: 詳細情報
    if [[ "${line}" =~ ^[[:space:]]+(.+\.go:[[:digit:]]+:.+) ]]; then
        echo  "  ↓ ${BASH_REMATCH[1]}"
        continue
    fi

    # INFO: ending行
    if [[ "${line}" =~ ^(FAIL|PASS)$ ]]; then
        continue
    fi
    if [[ "${line}" =~ ^(FAIL|ok)[[:space:]]+ ]]; then
        continue
    fi

done <<< "${output}"

end_ms=$(date +%s%3N)
elapsed_ms=$((end_ms - start_ms))

echo


# ------------------------------------------------------------
# Result
# ------------------------------------------------------------

# INFO: 時刻表示整形関数
format_seconds() {
    local ms=$1
    printf "%d.%02ds" $((ms / 1000)) $(((ms % 1000) / 10))
}

# INFO: NG率計算関数
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


echo "=== Result ==="
echo

if (( ng == 0 && status == 0 )); then
    printf "Total      : ${GREEN}●${RESET} PASS\n"
else
    printf "Total      : ${RED}●${RESET} FAIL\n"
fi

printf "Error Rate : %d tests / %d NG (%s)   %s\n" \
    "${total}" \
    "${ng}" \
    "$(ng_rate "${total}" "${ng}")" \
    "$(format_seconds "${elapsed_ms}")"

# ------------------------------------------------------------
# ending
# ------------------------------------------------------------
exit ${status}
