#!/usr/bin/env bash

set -e
cd "$(dirname "$0")/../app"

export ARIADNE_DATA_DIR="./devdata"

wails3 dev
