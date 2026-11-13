#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
./scripts/setup-security.sh
trap './scripts/compose.sh -p his-db-tests -f deploy/database-tests.compose.yaml down -v' EXIT HUP INT TERM
./scripts/compose.sh -p his-db-tests -f deploy/database-tests.compose.yaml up -d --wait
for service in identity staff catalog clinical laboratory lis; do
    port=55433
    case "$service" in
        identity) password=identity_local ;;
        staff) password=staff_local ;;
        catalog) password=catalog_local ;;
        clinical) password=clinical_local ;;
        laboratory) IFS= read -r password < deploy/security/services/database-laboratory.password ;;
        lis) IFS= read -r password < deploy/security/services/database-lis.password; port=55435 ;;
    esac
    (
        cd "services/$service"
        TEST_DATABASE_URL="postgres://${service}_user:${password}@127.0.0.1:${port}/${service}?sslmode=disable" \
            go test -race -count=1 ./...
    )
done
