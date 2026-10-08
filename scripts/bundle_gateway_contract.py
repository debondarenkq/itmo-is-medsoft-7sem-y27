#!/usr/bin/env python3
"""Build the public gateway specification from generated service contracts."""
import json
from pathlib import Path

root = Path(__file__).resolve().parent.parent
spec = {
    'openapi': '3.0.3',
    'info': {'title': 'HIS Gateway', 'version': '1.0.0'},
    'servers': [{'url': 'http://localhost:8080'}],
    'paths': {},
    'components': {'schemas': {}},
}
for service in ('staff', 'catalog', 'clinical'):
    source = json.loads((root / 'services' / service / 'api' / 'openapi.json').read_text())
    for path, value in source['paths'].items():
        if path in spec['paths']:
            raise SystemExit(f'Duplicate route: {path}')
        spec['paths'][path] = value
    for name, value in source['components']['schemas'].items():
        existing = spec['components']['schemas'].get(name)
        if existing is not None and existing != value:
            raise SystemExit(f'Conflicting schema: {name}')
        spec['components']['schemas'][name] = value
(root / 'desktop' / 'api' / 'openapi.json').write_text(json.dumps(spec, ensure_ascii=False, indent=2) + '\n')
