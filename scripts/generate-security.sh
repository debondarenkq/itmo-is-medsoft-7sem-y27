#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
isu_id=${1:-367907}
case "$isu_id" in ''|*[!0-9]*) printf '%s\n' 'ISU ID must contain digits' >&2; exit 1 ;; esac
security_dir=deploy/security
ca_dir="$security_dir/tls/ca"
server_dir="$security_dir/tls/server"
jwt_dir="$security_dir/jwt"
services_dir="$security_dir/services"
for name in ca.crt ca.key server.crt server.key signing.key service.key \
    tls/ca/ca.crt tls/ca/ca.key tls/server/server.crt tls/server/server.key jwt/signing.key services/service.key; do
    if [ -e "$security_dir/$name" ]; then
        printf '%s\n' 'Security files already exist; refusing to replace them. Run make certs to reuse or migrate the existing files.' >&2
        exit 1
    fi
done
umask 077
mkdir -p "$ca_dir" "$server_dir" "$jwt_dir" "$services_dir"
cat > "$ca_dir/ca.cnf" <<EOF
[req]
distinguished_name=dn
x509_extensions=ca
prompt=no
[dn]
O=ITMO
OU=MedSoft 2026/27 $isu_id
CN=HIS Teaching Root CA
[ca]
basicConstraints=critical,CA:TRUE,pathlen:0
keyUsage=critical,keyCertSign,cRLSign
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always
EOF
cat > "$server_dir/server.cnf" <<EOF
[req]
distinguished_name=dn
prompt=no
[dn]
O=ITMO
OU=MedSoft 2026/27 $isu_id
CN=localhost
[server]
basicConstraints=critical,CA:FALSE
keyUsage=critical,digitalSignature,keyEncipherment
extendedKeyUsage=serverAuth
subjectAltName=DNS:localhost,IP:127.0.0.1,IP:::1
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
EOF
openssl req -x509 -newkey rsa:3072 -nodes -days 3650 -sha256 \
    -config "$ca_dir/ca.cnf" -keyout "$ca_dir/ca.key" -out "$ca_dir/ca.crt"
openssl req -new -newkey rsa:3072 -nodes -config "$server_dir/server.cnf" \
    -keyout "$server_dir/server.key" -out "$server_dir/server.csr"
openssl x509 -req -in "$server_dir/server.csr" -CA "$ca_dir/ca.crt" \
    -CAkey "$ca_dir/ca.key" -CAcreateserial -days 365 -sha256 \
    -extfile "$server_dir/server.cnf" -extensions server -out "$server_dir/server.crt"
openssl genpkey -algorithm ED25519 -out "$jwt_dir/signing.key"
openssl rand -hex 32 > "$services_dir/service.key"
chmod 644 "$ca_dir/ca.crt" "$server_dir/server.crt"
printf '%s\n' 'Certificates and service keys created in deploy/security.'
