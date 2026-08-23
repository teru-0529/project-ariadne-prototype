#!/usr/bin/env bash

set -e
cd "$(dirname "$0")/../app"

wails3 build
