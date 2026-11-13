#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
./scripts/setup-security.sh
project=his-lis-tests
trap './scripts/compose.sh -p "$project" -f compose.yaml -f deploy/his-lis-tests.compose.yaml down -v' EXIT HUP INT TERM
./scripts/compose.sh -p "$project" -f compose.yaml -f deploy/his-lis-tests.compose.yaml up --no-build --pull never -d --wait
python3 scripts/check-his-lis.py
