#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Backend Test Runner
# ============================================================

set -u

GREEN='\E[32m\E[2m'
GRAY='\E[2m'
RESET='\E[0m'

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

    # INFO: Build Errorの可能性を通知
    if [[ "${line}" =~ ^#[[:space:]]+ ]]; then
        echo " 🔴 Potential Build Error"
        break
    fi

    # INFO: Test開始行削除
    if [[ "${line}" =~ ^===\ RUN[[:space:]]+ ]]; then
        continue
    fi

    # INFO: SKIP行削除
    if [[ "${line}" =~ ^---\ SKIP:[[:space:]]+ ]]; then
        continue
    fi

    # INFO: test結果行
    if [[ "${line}" =~ ^---\ (FAIL:|PASS:)[[:space:]]+Test_([^[:space:]]+)[[:space:]]+\((.+)\) ]]; then

        if [[ "${BASH_REMATCH[1]}" == "FAIL:" ]]; then
            printf " 🔴 ${BASH_REMATCH[2]//_/ } ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
            ((ng+=1))
        else
            printf " 🟢 ${BASH_REMATCH[2]//_/ } ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
        fi
        ((total+=1))
        continue
    fi

    # INFO: 詳細情報
    if [[ "${line}" =~ ^[[:space:]]+(.+\.go:[[:digit:]]+:.+) ]]; then
        printf "  ${GRAY}⇂${RESET} ${GRAY}%s${RESET}\n" "${BASH_REMATCH[1]}"
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

echo
echo "=== Result ==="
echo

if (( ng == 0 && status == 0 )); then
    printf "${GRAY}Total${RESET}      : 🟢 PASS\n"
else
    printf "${GRAY}Total${RESET}      : 🔴 FAIL\n"
fi

printf "${GRAY}Error Rate${RESET} : %d tests / %d NG (%s) ${GREEN}%s${RESET}\n" \
    "${total}" \
    "${ng}" \
    "$(ng_rate "${total}" "${ng}")" \
    "$(format_seconds "${elapsed_ms}")"

# ------------------------------------------------------------
# ending
# ------------------------------------------------------------
exit ${status}
