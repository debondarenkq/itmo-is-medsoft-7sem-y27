#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
for service in staff catalog clinical; do
    cp "contracts/$service.yaml" "services/$service/api/openapi.yaml"
    (cd "services/$service" && go generate ./api)
done
# Vendor external API descriptions so clinical-service still builds independently.
for service in staff catalog; do
    cp "contracts/$service.yaml" "services/clinical/internal/pkg/adapters/${service}api/openapi.yaml"
done
(cd services/clinical && go generate ./internal/pkg/adapters/staffapi ./internal/pkg/adapters/catalogapi)
python3 scripts/bundle_gateway_contract.py
docker run --rm -v "$PWD:/local" \
    openapitools/openapi-generator-cli:v7.26.0 generate \
    -c /local/macos/api/generator.yaml
python3 - <<'PY'
from pathlib import Path

for path in Path('macos/Generated').rglob('*.swift'):
    lines = (line.rstrip() for line in path.read_text().splitlines())
    path.write_text('\n'.join(lines).rstrip() + '\n')
PY
