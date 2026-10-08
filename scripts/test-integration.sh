#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
./scripts/compose.sh -p his-db-tests -f deploy/postgres-test.compose.yaml up -d --wait
trap './scripts/compose.sh -p his-db-tests -f deploy/postgres-test.compose.yaml down' EXIT
for service in staff catalog clinical; do
    case "$service" in
        staff) password=staff_local ;;
        catalog) password=catalog_local ;;
        clinical) password=clinical_local ;;
    esac
    (
        cd "services/$service"
        TEST_DATABASE_URL="postgres://${service}_user:${password}@127.0.0.1:55433/${service}?sslmode=disable" \
            go test -race -count=1 ./...
    )
done
