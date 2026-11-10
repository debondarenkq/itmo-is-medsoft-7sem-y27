#!/bin/sh
# Reuse existing teaching keys; never replace an identity during a normal launch.
set -eu
cd "$(dirname "$0")/.."

security_dir=deploy/security
umask 077
complete=true
exists=false
for name in tls/ca/ca.crt tls/ca/ca.key tls/server/server.crt tls/server/server.key jwt/signing.key services/service.key; do
    if [ -f "deploy/security/$name" ]; then
        exists=true
    else
        complete=false
    fi
done
if [ "$complete" = false ]; then
    if [ "$exists" = true ]; then
        printf '%s\n' 'Набор deploy/security неполный. Восстановите недостающие файлы; существующие ключи не заменены.' >&2
        exit 1
    fi
    ./scripts/generate-security.sh 367907
fi
openssl verify -CAfile deploy/security/tls/ca/ca.crt -purpose sslserver deploy/security/tls/server/server.crt
if ! openssl x509 -in deploy/security/tls/server/server.crt -checkend 0 -noout; then
    printf '%s\n' 'Серверный сертификат истёк. Требуется перевыпуск; обычный запуск не меняет ключи автоматически.' >&2
    exit 1
fi
printf '%s\n' 'Учебные сертификаты и ключи готовы. При повторном запуске используются те же файлы.'
