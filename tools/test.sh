#!/usr/bin/env bash

# ============================================================
# Project ARIADNE Prototype - Test Runner
# ============================================================

set -u

GRAY='\E[2m'
RESET='\E[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

backend_status=0
frontend_status=0

"${SCRIPT_DIR}/backend-test.sh" || backend_status=$?

echo

"${SCRIPT_DIR}/frontend-test.sh" || frontend_status=$?

echo
echo "=== Test Result ==="
echo

if [[ ${backend_status} -eq 0 ]]; then
    backend_result="🟢 PASS"
else
    backend_result="🔴 FAIL"
fi

if [[ ${frontend_status} -eq 0 ]]; then
    frontend_result="🟢 PASS"
else
    frontend_result="🔴 FAIL"
fi

if [[ ${backend_status} -eq 0 && ${frontend_status} -eq 0 ]]; then
    total_result="🟢 PASS"
    total_status=0
else
    total_result="🔴 FAIL"
    total_status=1
fi

printf "${GRAY}Backend${RESET}    : %s\n" "${backend_result}"
printf "${GRAY}Frontend${RESET}   : %s\n" "${frontend_result}"
echo "----------------"
printf "${GRAY}Total${RESET}      : %s\n" "${total_result}"

exit "${total_status}"
