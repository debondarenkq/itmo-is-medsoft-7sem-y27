#!/bin/sh
set -eu
IFS= read -r LABORATORY_DATABASE_PASSWORD < /security/database-laboratory.password
export LABORATORY_DATABASE_PASSWORD
psql -X -U "$POSTGRES_USER" -d postgres -v ON_ERROR_STOP=1 <<'SQL'
\getenv laboratory_password LABORATORY_DATABASE_PASSWORD
CREATE ROLE laboratory_user LOGIN PASSWORD :'laboratory_password';
CREATE DATABASE laboratory OWNER laboratory_user;
REVOKE CONNECT ON DATABASE laboratory FROM PUBLIC;
GRANT CONNECT ON DATABASE laboratory TO laboratory_user;
SQL
