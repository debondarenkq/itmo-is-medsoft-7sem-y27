#!/usr/bin/env python3
import argparse
import getpass
import json
import os
from pathlib import Path
import ssl
import sys
import tempfile
import urllib.error
import urllib.parse
import urllib.request

ROOT = Path(__file__).resolve().parent.parent
PAGE_SIZE = 200
REQUEST_TIMEOUT = 15


class NoRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, request, response, code, message, headers, new_url):
        return None


class API:
    def __init__(self, url, ca):
        address = urllib.parse.urlsplit(url)
        if (address.scheme != 'https' or not address.hostname or address.username
                or address.password or address.query or address.fragment):
            raise ValueError('--url должен содержать HTTPS-адрес gateway без логина и пароля.')
        self.url = url.rstrip('/')
        self.token = None
        self.refresh_token = None
        context = ssl.create_default_context(cafile=str(ca))
        context.minimum_version = ssl.TLSVersion.TLSv1_3
        self.opener = urllib.request.build_opener(
            urllib.request.HTTPSHandler(context=context), NoRedirects())
        self.operations = {}
        for description in ('/docs/security/openapi.json', '/docs/staff/openapi.json'):
            spec = self.request('GET', description)
            for path, methods in spec['paths'].items():
                for method, operation in methods.items():
                    if isinstance(operation, dict) and 'operationId' in operation:
                        self.operations[operation['operationId']] = method.upper(), path

    def request(self, method, path, body=None, parameters=None):
        if parameters:
            path += "?" + urllib.parse.urlencode(parameters)
        payload = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request(self.url + path, data=payload, method=method)
        if payload is not None:
            request.add_header('Content-Type', 'application/json')
        if self.token:
            request.add_header('Authorization', 'Bearer ' + self.token)
        try:
            with self.opener.open(request, timeout=REQUEST_TIMEOUT) as response:
                content = response.read()
        except urllib.error.HTTPError as error:
            raise RuntimeError(f'{method} {path}: HTTP {error.code}') from None
        except urllib.error.URLError as error:
            raise RuntimeError(f'Gateway недоступен: {error.reason}') from None
        return json.loads(content) if content else None

    def call(self, operation, body=None, parameters=None):
        method, path = self.operations[operation]
        return self.request(method, path, body, parameters)

    def pages(self, operation, **parameters):
        items = []
        while True:
            page = self.call(operation, parameters={**parameters, 'limit': PAGE_SIZE, 'offset': len(items)})
            items.extend(page['items'])
            if len(items) >= page['total']:
                return items
            if not page['items']:
                raise RuntimeError('Gateway вернул пустую страницу до конца списка.')

    def login(self):
        if not self.call('GetSetup')['configured']:
            raise RuntimeError('Сначала настройте пароль и TOTP администратора в HISAdmin.')
        password = getpass.getpass('Пароль администратора: ')
        totp = getpass.getpass('Текущий TOTP-код: ')
        session = self.call('Login',
                            {'username': 'admin', 'password': password, 'totp': totp})
        self.token = session['access_token']
        self.refresh_token = session['refresh_token']

    def logout(self):
        if self.refresh_token:
            self.call('Logout', {'refresh_token': self.refresh_token})
            self.token = self.refresh_token = None


def save_credentials(path, credentials):
    path.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.NamedTemporaryFile(dir=path.parent, delete=False) as stream:
        temporary = Path(stream.name)
        try:
            os.chmod(temporary, 0o600)
            stream.write((json.dumps(list(credentials.values()), ensure_ascii=False, indent=2) + '\n').encode())
            stream.flush()
            os.replace(temporary, path)
        finally:
            temporary.unlink(missing_ok=True)


def create_accounts(api, profiles, output):
    staff = api.pages('ListStaff', include_deleted='true')
    accounts = {account['username']: account for account in api.pages('ListAccounts')}
    previous = json.loads(output.read_text()) if output.exists() else []
    credentials = {item['username']: item for item in previous}
    created = 0
    for profile in profiles:
        if profile.get('deleted') or 'username' not in profile:
            continue
        username = profile['username']
        fields = ('first_name', 'last_name', 'position')
        matching = [worker for worker in staff
                    if all(worker[field] == profile[field] for field in fields)]
        if len(matching) > 1:
            raise RuntimeError(f'{username}: найдено несколько одинаковых сотрудников; выберите сотрудника вручную.')
        account = accounts.get(username)
        if account:
            if (not matching or matching[0]['deleted_at'] is not None
                    or account.get('staff_id') != matching[0]['id']):
                raise RuntimeError(f'{username}: логин занят другим или удалённым сотрудником.')
            saved = credentials.get(username, {})
            if saved.get('account_id') != account['id'] or saved.get('staff_id') != account['staff_id']:
                credentials.pop(username, None)
            print(f'{username}: уже существует; пароль не изменён.')
            continue
        if matching:
            worker = matching[0]
            if worker['deleted_at'] is not None:
                raise RuntimeError(f'{username}: сотрудник удалён; восстановите данные вручную.')
        else:
            worker = api.call('CreateStaff',
                              {field: profile[field] for field in (*fields, 'role')})
            staff.append(worker)
        result = api.call('CreateAccount', {'username': username, 'staff_id': worker['id']})
        account = result['account']
        accounts[username] = account
        credentials[username] = dict(username=username, role=account['role'],
                                     account_id=account['id'], staff_id=worker['id'],
                                     password=result['password'])
        save_credentials(output, credentials)
        created += 1
        print(f'{username}: создан.')
    save_credentials(output, credentials)
    return created


def main():
    parser = argparse.ArgumentParser(description='Create demo staff and accounts through the public HTTPS API.')
    parser.add_argument('--url', default='https://localhost:8443', help='HTTPS-адрес gateway')
    parser.add_argument('--ca', type=Path, default=ROOT / 'deploy/security/tls/ca/ca.crt')
    parser.add_argument('--output', type=Path, default=ROOT / 'build/demo-accounts.json')
    args = parser.parse_args()
    profiles = json.loads((ROOT / 'examples/demo-staff.json').read_text())
    api = API(args.url, args.ca)
    api.login()
    try:
        created = create_accounts(api, profiles, args.output)
    finally:
        try:
            api.logout()
        except RuntimeError as error:
            print(f'Не удалось завершить сеанс скрипта: {error}', file=sys.stderr)
    print(f'Создано аккаунтов: {created}. Выданные пароли сохранены в {args.output} (права 0600).')
    print('Пароли существующих аккаунтов при необходимости сбрасываются через HISAdmin.')


if __name__ == '__main__':
    try:
        main()
    except (RuntimeError, ValueError, OSError) as error:
        sys.exit(str(error))
