#!/usr/bin/env python3
"""Build the public gateway specification from generated service contracts."""
import json
from pathlib import Path

root = Path(__file__).resolve().parent.parent
spec = {
    'openapi': '3.0.3',
    'info': {'title': 'HIS Gateway', 'version': '1.0.0'},
    'servers': [{'url': 'https://localhost:8443'}],
    'paths': {},
    'components': {'schemas': {}, 'securitySchemes': {'BearerAuth': {'type': 'http', 'scheme': 'bearer', 'bearerFormat': 'JWT'}}},
    'security': [{'BearerAuth': []}],
}
for service in ('staff', 'catalog', 'clinical', 'staff/api/security'):
    source_path = root / 'services' / service / ('openapi.json' if service == 'staff/api/security' else 'api/openapi.json')
    source = json.loads(source_path.read_text())
    for path, value in source['paths'].items():
        if path in spec['paths']:
            raise SystemExit(f'Duplicate route: {path}')
        spec['paths'][path] = value
    for name, value in source['components']['schemas'].items():
        existing = spec['components']['schemas'].get(name)
        if existing is not None and existing != value:
            raise SystemExit(f'Conflicting schema: {name}')
        spec['components']['schemas'][name] = value
(root / 'macos' / 'api' / 'openapi.json').write_text(json.dumps(spec, ensure_ascii=False, indent=2) + '\n')
