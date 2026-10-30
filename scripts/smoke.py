#!/usr/bin/env python3
"""Check lab 2 through a TLS-verified gateway using synthetic data.

Use --bootstrap only with a fresh, disposable Compose project. Without it the
script prompts for the existing administrator password and current TOTP code.
Passwords and tokens are kept in memory and never printed.
"""
import argparse
import base64
import getpass
import hashlib
import hmac
import json
from pathlib import Path
import secrets
import ssl
import socket
import struct
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid


def totp(secret, timestamp=None):
    step = int(time.time() if timestamp is None else timestamp) // 30
    digest = hmac.new(base64.b32decode(secret), struct.pack('>Q', step), hashlib.sha1).digest()
    offset = digest[-1] & 15
    return f'{(struct.unpack(">I", digest[offset:offset + 4])[0] & 0x7fffffff) % 1000000:06d}'


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--url', default='https://localhost:8443')
    parser.add_argument('--ca', default=str(Path(__file__).resolve().parent.parent / 'deploy/security/tls/ca/ca.crt'))
    parser.add_argument('--bootstrap', action='store_true')
    args = parser.parse_args()
    assert urllib.parse.urlparse(args.url).scheme == 'https', 'HTTPS is required'
    context = ssl.create_default_context(cafile=args.ca)
    context.minimum_version = ssl.TLSVersion.TLSv1_3
    target = urllib.parse.urlparse(args.url)
    with socket.create_connection((target.hostname, target.port or 443), timeout=5) as connection:
        with context.wrap_socket(connection, server_hostname=target.hostname) as secured:
            assert secured.version() == 'TLSv1.3'
            issuer = dict(part[0] for part in secured.getpeercert()['issuer'])
            assert issuer['organizationName'] == 'ITMO'
            assert issuer['organizationalUnitName'] == 'MedSoft 2026/27 367907'
    legacy = ssl.create_default_context(cafile=args.ca)
    legacy.maximum_version = ssl.TLSVersion.TLSv1_2
    try:
        with socket.create_connection((target.hostname, target.port or 443), timeout=5) as connection:
            with legacy.wrap_socket(connection, server_hostname=target.hostname):
                raise AssertionError('The gateway accepted TLS 1.2')
    except ssl.SSLError:
        pass
    token = None

    def call(method, path, body=None, status=200, bearer=None):
        payload = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request(args.url.rstrip('/') + path, data=payload, method=method)
        if payload is not None:
            request.add_header('Content-Type', 'application/json')
        access = token if bearer is None else bearer
        if access:
            request.add_header('Authorization', 'Bearer ' + access)
        try:
            response = urllib.request.urlopen(request, context=context, timeout=10)
        except urllib.error.HTTPError as error:
            response = error
        with response:
            actual = response.status
            data = response.read()
        # Do not include response bodies in failures: they may contain credentials.
        assert actual == status, f'{method} {path}: expected HTTP {status}, got {actual}'
        return json.loads(data) if data else None

    for path in ['/api/v1/staff', '/api/v1/diagnoses', '/api/v1/patients', '/api/v1/accounts', '/api/v1/audit']:
        call('GET', path, status=401)
    setup = call('GET', '/api/v1/auth/setup')
    if args.bootstrap:
        assert not setup['configured'], '--bootstrap requires a fresh database'
        password = "Aa1" + secrets.token_urlsafe(24)
        secret = setup['secret']
        call('POST', '/api/v1/auth/setup', {'password': password, 'confirmation': password + 'x', 'secret': secret, 'totp': totp(secret)}, 400)
        call('POST', '/api/v1/auth/setup', {'password': password, 'confirmation': password, 'secret': secret, 'totp': 'invalid'}, 401)
        assert not call('GET', '/api/v1/auth/setup')['configured']
        admin = call('POST', '/api/v1/auth/setup', {'password': password, 'confirmation': password, 'secret': secret, 'totp': totp(secret)})
        assert call('GET', '/api/v1/auth/setup')['configured']
        # The TOTP step accepted during setup cannot be reused for login.
        call('POST', '/api/v1/auth/login', {'username': 'admin', 'password': password, 'totp': totp(secret)}, 401)
        admin = call('POST', '/api/v1/auth/login', {'username': 'admin', 'password': password, 'totp': totp(secret, time.time() + 30)})
    else:
        assert setup['configured'], 'Set up the administrator in the app or use --bootstrap on an isolated project'
        admin = call('POST', '/api/v1/auth/login', {'username': 'admin', 'password': getpass.getpass('Administrator password: '), 'totp': input('TOTP: ').strip()})
    call('POST', '/api/v1/auth/login', {'username': 'missing-user', 'password': 'wrong-password'}, 401)
    previous = admin['refresh_token']
    admin = call('POST', '/api/v1/auth/refresh', {'refresh_token': previous})
    assert admin['refresh_token'] != previous
    call('POST', '/api/v1/auth/refresh', {'refresh_token': previous}, 401)
    token = admin['access_token']
    call('GET', '/api/v1/patients', status=403)
    staff_input = {'first_name': 'А', 'last_name': 'Б', 'position': 'Медработник'}
    workers = {}
    passwords = {}
    for role in ['doctor', 'nurse', 'laboratory-technician', 'other']:
        staff = call('POST', '/api/v1/staff', {**staff_input, 'role': role}, 201)
        account = call('POST', '/api/v1/accounts', {'staff_id': staff['id'], 'username': 'u' + uuid.uuid4().hex[:12]}, 201)
        generated = account['password']
        passwords[role] = generated
        assert len(generated) == 12 and all(c in 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789' for c in generated)
        session = call('POST', '/api/v1/auth/login', {'username': account['account']['username'], 'password': generated})
        workers[role] = (staff, account['account'], session)
    doctor, doctor_account, doctor_session = workers['doctor']
    account_search = urllib.parse.urlencode({'search': doctor_account['username'], 'is_active': 'true'})
    matching = call('GET', '/api/v1/accounts?' + account_search)
    assert matching['total'] == 1 and matching['items'][0]['staff_name'] == 'Б А'
    assert matching['items'][0]['role'] == 'doctor'
    new_username = 'u' + uuid.uuid4().hex[:12]
    doctor_account.update(call('PUT', f'/api/v1/accounts/{doctor_account["id"]}', {'is_active': True, 'username': new_username}))
    call('POST', '/api/v1/accounts', {'staff_id': doctor['id'], 'username': 'duplicate'}, 409)
    call('POST', '/api/v1/staff', {**staff_input, 'role': 'administrator'}, 400)
    file_data = json.loads((Path(__file__).resolve().parent.parent / 'examples/diagnoses.json').read_text())
    imported = call('POST', '/api/v1/diagnoses/import', file_data)['items']
    again = call('POST', '/api/v1/diagnoses/import', file_data)['items']
    assert [d['id'] for d in imported] == [d['id'] for d in again]
    did = imported[0]['id']
    patient_input = {'first_name': 'Иван', 'last_name': 'Тестовый', 'middle_name': None, 'birth_date': '2000-01-01', 'administrative_sex': 'UNKNOWN', 'comment': ''}
    token = doctor_session['access_token']
    patient = call('POST', '/api/v1/patients', patient_input, 201)
    pid = patient['id']
    disposable = call('POST', '/api/v1/patients', patient_input, 201)
    call('DELETE', f'/api/v1/patients/{disposable["id"]}', status=204)
    call('PUT', f'/api/v1/patients/{pid}', {**patient_input, 'comment': 'Учебный пациент'})
    call('POST', f'/api/v1/patients/{pid}/record', {'staff_id': workers['nurse'][0]['id']}, 400)
    record = call('POST', f'/api/v1/patients/{pid}/record', {'staff_id': doctor['id']}, 201)
    assert record['number'] > 0
    patient_card = call('GET', f'/api/v1/patients/{pid}')
    assert patient_card['record_id'] == record['id'] and patient_card['record_number'] == record['number']
    number_search = urllib.parse.urlencode({'search': f"ЭМК-{record['number']:06d}"})
    assert call('GET', '/api/v1/patients?' + number_search)['total'] == 1
    record_search = urllib.parse.urlencode({'search': record['id'][:8], 'has_record': 'true'})
    found = call('GET', '/api/v1/patients?' + record_search)
    assert found['total'] == 1 and found['items'][0]['id'] == pid
    rid = record['id']
    path = f'/api/v1/records/{rid}'
    call('DELETE', f'/api/v1/patients/{pid}', status=409)

    def save(version, commands, status=200):
        return call('POST', path + '/changes', {'staff_id': doctor['id'], 'expected_version': version, 'commands': commands}, status)

    save(1, [{'type': 'add_prescription', 'text': 'Пить воду'}, {'type': 'remove_diagnosis', 'id': str(uuid.uuid4())}], 404)
    assert call('GET', path)['version'] == 1
    added = save(1, [{'type': 'add_diagnosis', 'diagnosis_id': did}, {'type': 'add_prescription', 'text': 'Пить больше воды'}])
    rxid = added['state']['prescriptions'][0]['id']
    entry_id = added['state']['diagnoses'][0]['id']
    save(2, [{'type': 'add_diagnosis', 'diagnosis_id': did}], 409)
    save(1, [{'type': 'add_prescription', 'text': 'Назначение'}], 409)
    edited = save(2, [{'type': 'edit_prescription', 'id': rxid, 'text': 'Пить воду после еды'}])
    assert edited['state']['prescriptions'][0]['id'] == rxid
    cancelled = save(3, [{'type': 'cancel_prescription', 'id': rxid}, {'type': 'remove_diagnosis', 'id': entry_id}])
    assert cancelled['state']['prescriptions'][0]['status'] == 'cancelled'
    for snapshot in [added, edited, cancelled]:
        assert call('GET', path + f'/state?version={snapshot["version"]}')['state'] == snapshot['state']
        query = urllib.parse.urlencode({'at': snapshot['updated_at']})
        assert call('GET', path + '/state?' + query)['state'] == snapshot['state']
    assert call('GET', path + '/history')['total'] == 6
    for role, (staff, account, session) in workers.items():
        token = session['access_token']
        call('GET', path)
        call('GET', '/api/v1/accounts', status=403)
        call('GET', '/api/v1/audit', status=403)
        if role != 'doctor':
            call('POST', path + '/changes', {'staff_id': staff['id'], 'expected_version': 4, 'commands': []}, 403)
        if role in ['other', 'laboratory-technician']:
            call('POST', '/api/v1/patients', patient_input, 403)
            call('POST', f'/api/v1/patients/{pid}/record', {'staff_id': staff['id']}, 403)
        if role == 'nurse':
            nurse_patient = call('POST', '/api/v1/patients', patient_input, 201)
            call('POST', f'/api/v1/patients/{nurse_patient["id"]}/record', {'staff_id': staff['id']}, 201)
    token = admin['access_token']
    # Deactivation, password changes, role changes and deletion each revoke refresh.
    other_staff, other_account, other_session = workers['other']
    call('PUT', f'/api/v1/accounts/{other_account["id"]}', {'is_active': False})
    call('POST', '/api/v1/auth/login', {'username': other_account['username'], 'password': passwords['other']}, 401)
    call('POST', '/api/v1/auth/refresh', {'refresh_token': other_session['refresh_token']}, 401)
    call('PUT', f'/api/v1/accounts/{other_account["id"]}', {'is_active': True})
    reset = call('POST', f'/api/v1/accounts/{doctor_account["id"]}/password')
    call('POST', '/api/v1/auth/login', {'username': doctor_account['username'], 'password': passwords['doctor']}, 401)
    call('POST', '/api/v1/auth/refresh', {'refresh_token': doctor_session['refresh_token']}, 401)
    session = call('POST', '/api/v1/auth/login', {'username': doctor_account['username'], 'password': reset['password']})
    call('PUT', f'/api/v1/accounts/{doctor_account["id"]}', {'is_active': True, 'role': 'nurse'})
    call('POST', '/api/v1/auth/refresh', {'refresh_token': session['refresh_token']}, 401)
    changed = call('POST', '/api/v1/auth/login', {'username': doctor_account['username'], 'password': reset['password']})
    assert changed['role'] == 'nurse'
    call('POST', path + '/changes', {'staff_id': doctor['id'], 'expected_version': 4, 'commands': []}, 403, bearer=changed['access_token'])
    staff, account, session = workers['laboratory-technician']
    call('DELETE', f'/api/v1/staff/{staff["id"]}', status=204)
    call('POST', '/api/v1/auth/refresh', {'refresh_token': session['refresh_token']}, 401)
    call('PUT', f'/api/v1/accounts/{account["id"]}', {'is_active': True}, 409)
    accounts = call('GET', '/api/v1/accounts?limit=200')['items']
    assert not next(a for a in accounts if a['id'] == account['id'])['is_active']
    assert all('password_hash' not in a for a in accounts)
    events = call('GET', '/api/v1/audit?limit=200')['items']
    assert events and any(e['actor'].startswith('anonymous:') for e in events)
    assert any('отозвано refresh токенов: 1' in e['message'] for e in events)
    assert events == sorted(events, key=lambda e: (e['event_time'], e['event_id']))
    filtered = call('GET', '/api/v1/audit?event_type=patient_created&limit=200')
    assert filtered['total'] > 0 and all(e['event_type'] == 'patient_created' for e in filtered['items'])
    assert all(e['actor_name'] for e in filtered['items'])
    future = urllib.parse.urlencode({'from': '3000-01-01T00:00:00Z'})
    assert call('GET', '/api/v1/audit?' + future)['total'] == 0
    invalid_period = urllib.parse.urlencode({'from': '3000-01-01T00:00:00Z', 'to': '2000-01-01T00:00:00Z'})
    call('GET', '/api/v1/audit?' + invalid_period, status=400)
    encoded = json.dumps(events)
    assert reset['password'] not in encoded and admin['refresh_token'] not in encoded
    call('POST', '/api/v1/auth/logout', {'refresh_token': admin['refresh_token']})
    call('POST', '/api/v1/auth/refresh', {'refresh_token': admin['refresh_token']}, 401)
    print('PASS: TLS 1.3, admin setup/TOTP, login, JWT, rotation/replay, all roles, actor binding, accounts, revocation, audit, EMR/history/rollback')


if __name__ == '__main__':
    main()
