#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Frontend Test Runner
# ============================================================

set -u

GREEN='\E[32m\E[2m'
PINK='\E[1;37;48;2;255;153;180m'
RED_BOLD='\E[1;91m'
GRAY='\E[2m'
RESET='\E[0m'


# ------------------------------------------------------------
# Initial
# ------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
FRONTEND_DIR="${PROJECT_ROOT}/app/frontend"

status=0

# ------------------------------------------------------------
# Frontend
# ------------------------------------------------------------

echo
echo "=== Frontend Test ==="
echo

start_ms=$(date +%s%3N)

# INFO: npm run test 実行
output="$(
    cd "${FRONTEND_DIR}" &&
    npm run test -- --reporter=tree 2>&1
)"

status=$?

total=0
ng=0

while IFS= read -r line; do

    plain_line="$(printf '%s' "${line}" | sed $'s/\033\\[[0-9;]*m//g')"
    # printf '%q\n' "${line}"
    # printf 'PLAIN=[%s]\n' "${plain_line}"


    # INFO: testファイル行
    if [[ "${plain_line}" =~ ^[[:space:]](❯|✓)[[:space:]](.+)[[:space:]]([0-9]+ms)$ ]]; then

        if [[ "${BASH_REMATCH[1]}" == "❯" ]]; then
            printf "🔴 ${GRAY}${BASH_REMATCH[2]}${RESET} ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
        else
            printf "🟢 ${GRAY}${BASH_REMATCH[2]}${RESET} ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
        fi
        continue
    fi

    # INFO: test結果行
    if [[ "${plain_line}" =~ ^[[:space:]]{3}(×|✓)[[:space:]](.+)[[:space:]]([0-9]+ms)$ ]]; then

        if [[ "${BASH_REMATCH[1]}" == "×" ]]; then
            printf "  🔴 ${BASH_REMATCH[2]} ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
            ((ng+=1))
        else
            printf "  🟢 ${BASH_REMATCH[2]} ${GREEN}${BASH_REMATCH[3]}${RESET}\n"
        fi
        ((total+=1))
        continue
    fi

    # INFO: エラー有無判断
    if [[ "${plain_line}" =~ Failed[[:space:]]Tests[[:space:]]+[0-9]+ ]]; then
        echo
        echo "=== Errors ==="
        continue
    fi

    # INFO: Error詳細行
    if [[ "${plain_line}" =~ ^[[:space:]]FAIL[[:space:]](.+)[[:space:]]\>[[:space:]](.+) ]]; then
        echo
        printf "${PINK} FAIL ${RESET} ${BASH_REMATCH[1]} ${GRAY}>${RESET} ${BASH_REMATCH[2]}\n"
        continue
    fi

    # INFO: Error詳細行
    if [[ "${plain_line}" =~ ^(.+Error):[[:space:]](.+) ]]; then
        printf "${RED_BOLD}${BASH_REMATCH[1]}${RESET}: ${GRAY}${BASH_REMATCH[2]}${RESET} \n"
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
