#!/usr/bin/env python3
"""Read-only checks for the stack used at the lab defense; never prints secrets."""
import json
from pathlib import Path
import socket
import ssl
import subprocess
import sys
import urllib.error
import urllib.request

ROOT = Path(__file__).resolve().parent.parent


def main():
    failures = []

    def check(label, action):
        try:
            detail = action()
            print(f'OK: {label}' + (f' — {detail}' if detail else ''), flush=True)
        except Exception as error:
            failures.append(label)
            print(f'ОШИБКА: {label} — {error}', flush=True)

    def docker_status():
        result = subprocess.run(
            [str(ROOT / 'scripts/compose.sh'), 'ps', '--all', '--format', 'json'],
            cwd=ROOT, check=True, capture_output=True, text=True,
        )
        # Compose versions return either one JSON array or JSON objects per line.
        content = result.stdout.strip()
        rows = json.loads(content) if content.startswith('[') else [json.loads(line) for line in content.splitlines()]
        expected = {'postgres', 'staff-service', 'catalog-service', 'clinical-service', 'nginx'}
        ready = {row['Service'] for row in rows if row.get('State') == 'running' and row.get('Health') == 'healthy'}
        missing = expected - ready
        if missing:
            raise RuntimeError('не готовы: ' + ', '.join(sorted(missing)) + '; выполните make up')
        return 'все 5 компонентов готовы'

    def app_exists():
        if not (ROOT / 'macos/build/HIS.app/Contents/MacOS/HIS').is_file():
            raise RuntimeError('выполните make macos-build')
        return 'macos/build/HIS.app'

    def tls_connection():
        context = ssl.create_default_context(cafile=str(ROOT / 'deploy/security/tls/ca/ca.crt'))
        context.minimum_version = ssl.TLSVersion.TLSv1_3
        with socket.create_connection(('localhost', 8443), timeout=5) as connection:
            with context.wrap_socket(connection, server_hostname='localhost') as secured:
                if secured.version() != 'TLSv1.3':
                    raise RuntimeError('ожидался TLS 1.3')
        with urllib.request.urlopen('https://localhost:8443/live', context=context, timeout=5) as response:
            if response.status != 200:
                raise RuntimeError(f'проверка gateway: HTTP {response.status}')
        try:
            with urllib.request.urlopen('https://localhost:8443/api/v1/patients', context=context, timeout=5):
                raise RuntimeError('API пациентов разрешил доступ без входа')
        except urllib.error.HTTPError as error:
            if error.code != 401:
                raise RuntimeError(f'ожидался отказ 401, получен {error.code}') from error
        return 'проверены CA, имя localhost, TLS 1.3 и отказ API без авторизации'

    check('Docker Compose', docker_status)
    check('Нативное приложение', app_exists)
    check('Защищённое подключение', tls_connection)
    print('Доверие к CA именно в macOS проверяется в HIS → Сертификаты сервера → Проверить подключение.')
    if failures:
        return 1
    print('Готово к запуску: make demo')
    return 0


if __name__ == '__main__':
    sys.exit(main())
