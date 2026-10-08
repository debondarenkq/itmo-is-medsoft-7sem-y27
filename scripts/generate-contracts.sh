#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
for service in staff catalog clinical; do
    (cd "services/$service" && go generate ./api)
done
# Vendor external API descriptions so clinical-service still builds independently.
for service in staff catalog; do
    cp "services/$service/api/openapi.yaml" "services/clinical/internal/pkg/adapters/${service}api/openapi.yaml"
done
(cd services/clinical && go generate ./internal/pkg/adapters/staffapi ./internal/pkg/adapters/catalogapi)
python3 scripts/bundle_gateway_contract.py
docker run --rm -v "$PWD:/local" \
    openapitools/openapi-generator-cli:v7.26.0 generate \
    -c /local/desktop/api/generator.yaml
if command -v dart >/dev/null 2>&1; then
    dart format desktop/packages/his_api/lib
else
    docker run --rm -e CI=true -e PUB_CACHE=/pub-cache \
        -v his-flutter-pub-cache:/pub-cache \
        -v "$PWD/desktop:/app" -w /app \
        ghcr.io/cirruslabs/flutter@sha256:46691e311715845de03a3ba4753a475476936805b29431b1f00f1816981033f8 \
        dart format packages/his_api/lib
fi
