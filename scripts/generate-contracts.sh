#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
for service in staff catalog clinical laboratory; do
    cp "contracts/$service.yaml" "services/$service/api/openapi.yaml"
    if [ "$service" = laboratory ]; then cp contracts/fhir.yaml "services/$service/api/fhir.yaml"; fi
    (cd "services/$service" && go generate ./api)
done
cp contracts/audit.yaml services/identity/api/audit/openapi.yaml
cp contracts/audit.yaml services/identity/api/security/audit.yaml
(cd services/identity && go generate ./api/audit ./internal/models)
for service in staff catalog clinical laboratory; do
    cp contracts/audit.yaml "services/$service/internal/pkg/adapters/audit/openapi.yaml"
    (cd "services/$service" && go generate ./internal/pkg/adapters/audit)
done
cp contracts/security.yaml services/identity/api/security/openapi.yaml
(cd services/identity && go generate ./api/security && go run ./api/security/export > api/security/openapi.json)
cp contracts/internal-identity.yaml services/identity/api/internalapi/openapi.yaml
(cd services/identity/api/internalapi && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/internal-identity.yaml services/staff/api/internalapi/openapi.yaml
(cd services/staff/api/internalapi && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/internal-identity.yaml services/staff/internal/pkg/adapters/security/openapi.yaml
(cd services/staff/internal/pkg/adapters/security && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/internal-identity.yaml services/catalog/internal/pkg/adapters/security/openapi.yaml
(cd services/catalog/internal/pkg/adapters/security && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/internal-identity.yaml services/clinical/internal/pkg/adapters/security/openapi.yaml
(cd services/clinical/internal/pkg/adapters/security && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/internal-identity.yaml services/identity/internal/pkg/adapters/staff/openapi.yaml
(cd services/identity/internal/pkg/adapters/staff && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
for service in staff catalog; do
    cp "contracts/$service.yaml" "services/clinical/internal/pkg/adapters/${service}api/openapi.yaml"
done
(cd services/clinical && go generate ./internal/pkg/adapters/staffapi ./internal/pkg/adapters/catalogapi)
cp contracts/internal-identity.yaml services/laboratory/internal/pkg/adapters/security/openapi.yaml
(cd services/laboratory/internal/pkg/adapters/security && go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml)
cp contracts/clinical.yaml services/laboratory/internal/pkg/adapters/clinicalapi/openapi.yaml
(cd services/laboratory && go generate ./internal/pkg/adapters/clinicalapi)
cp contracts/lis.yaml services/lis/api/openapi.yaml
cp contracts/fhir.yaml services/lis/api/fhir.yaml
(cd services/lis && go generate ./api)
for service in lis identity; do
    cp contracts/laboratory.yaml "services/$service/internal/pkg/adapters/laboratoryapi/openapi.yaml"
    cp contracts/fhir.yaml "services/$service/internal/pkg/adapters/laboratoryapi/fhir.yaml"
    (cd "services/$service" && go generate ./internal/pkg/adapters/laboratoryapi)
done
cp contracts/security.yaml services/lis/internal/pkg/adapters/identityapi/openapi.yaml
cp contracts/audit.yaml services/lis/internal/pkg/adapters/identityapi/audit.yaml
(cd services/lis && go generate ./internal/pkg/adapters/identityapi)
python3 scripts/bundle_gateway_contract.py
./scripts/docker.sh run --rm -v "$PWD:/local" \
    openapitools/openapi-generator-cli:v7.26.0 generate \
    -c /local/macos/api/generator.yaml
python3 - <<'PY'
from pathlib import Path

for path in Path('macos/Generated').rglob('*.swift'):
    lines = (line.rstrip() for line in path.read_text().splitlines())
    path.write_text('\n'.join(lines).rstrip() + '\n')
PY
