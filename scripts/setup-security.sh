#!/bin/sh
# Reuse existing teaching keys; never replace an identity during a normal launch.
set -eu
cd "$(dirname "$0")/.."

security_dir=deploy/security
umask 077
# Move the former flat layout without replacing keys or changing their permissions.
legacy_files='ca.crt:tls/ca/ca.crt ca.key:tls/ca/ca.key ca.cnf:tls/ca/ca.cnf ca.srl:tls/ca/ca.srl
server.crt:tls/server/server.crt server.key:tls/server/server.key server.cnf:tls/server/server.cnf server.csr:tls/server/server.csr
signing.key:jwt/signing.key service.key:services/service.key'
for mapping in $legacy_files; do
    old_name=${mapping%%:*}
    new_name=${mapping#*:}
    if [ -e "$security_dir/$old_name" ] && [ -e "$security_dir/$new_name" ]; then
        printf 'Есть и старый, и новый файл %s. Перенос остановлен без перезаписи.\n' "$old_name" >&2
        exit 1
    fi
done
for mapping in $legacy_files; do
    old_name=${mapping%%:*}
    new_name=${mapping#*:}
    if [ -f "$security_dir/$old_name" ]; then
        mkdir -p "$security_dir/${new_name%/*}"
        mv "$security_dir/$old_name" "$security_dir/$new_name"
    fi
done

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
