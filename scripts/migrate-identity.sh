#!/bin/sh
# One-time offline transfer. Never resets users, TOTP secrets, tokens, or audit history.
set -eu
cd "$(dirname "$0")/.."
./scripts/compose.sh up --no-build --pull never -d --wait postgres
# init.sql is not rerun for existing volumes. Provision only the missing identity role/database.
./scripts/compose.sh exec -T postgres psql -X -v ON_ERROR_STOP=1 -U his_admin -d postgres <<'SQL'
SELECT 'CREATE ROLE identity_user LOGIN PASSWORD ''identity_local'''
 WHERE NOT EXISTS(SELECT FROM pg_roles WHERE rolname='identity_user')\gexec
SELECT 'CREATE DATABASE identity OWNER identity_user'
 WHERE NOT EXISTS(SELECT FROM pg_database WHERE datname='identity')\gexec
REVOKE CONNECT ON DATABASE identity FROM PUBLIC;
GRANT CONNECT ON DATABASE identity TO identity_user;
SQL
legacy=$(./scripts/compose.sh exec -T postgres psql -X -At -U his_admin -d staff -c "SELECT to_regclass('public.accounts') IS NOT NULL")
if [ "$legacy" = t ]; then
    printf '%s\n' 'Перенос аккаунтов и аудита в identity: временно останавливаем API.'
    ./scripts/compose.sh stop nginx clinical-service catalog-service staff-service identity-service
    mkdir -p build
    umask 077
    if [ ! -f build/staff-before-identity.dump ]; then
        ./scripts/compose.sh exec -T postgres pg_dump -U his_admin -d staff -Fc > build/staff-before-identity.dump.partial
        mv build/staff-before-identity.dump.partial build/staff-before-identity.dump
    fi
    ./scripts/compose.sh run --rm --no-deps --pull never \
        -e LEGACY_DATABASE_URL='postgres://staff_user:staff_local@postgres:5432/staff?sslmode=disable' \
        identity-service migrate
fi
