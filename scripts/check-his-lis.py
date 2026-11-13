#!/usr/bin/env python3
"""Verify HIS/LIS integration on the disposable his-lis-tests stack; never print credentials."""
import argparse
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime, timezone
import json
from pathlib import Path
import secrets
import ssl
import urllib.error
import urllib.request
import uuid

from smoke import totp


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--his', default='https://localhost:58443')
    parser.add_argument('--lis', default='https://localhost:58444')
    args = parser.parse_args()
    context = ssl.create_default_context(cafile=str(Path(__file__).resolve().parent.parent / 'deploy/security/tls/ca/ca.crt'))
    context.minimum_version = ssl.TLSVersion.TLSv1_3

    def call(method, path, body=None, token=None, status=200, lis=False, fhir=False):
        payload = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request((args.lis if lis else args.his) + path, data=payload, method=method)
        if payload is not None:
            request.add_header('Content-Type', 'application/fhir+json' if fhir else 'application/json')
        if token:
            request.add_header('Authorization', 'Bearer ' + token)
        try:
            response = urllib.request.urlopen(request, context=context, timeout=25)
        except urllib.error.HTTPError as error:
            response = error
        with response:
            data = response.read()
            diagnostic = ''
            if response.status != status and (path.startswith('/lis/v1/results') or path.startswith('/api/v1/fhir/')):
                problem = json.loads(data)
                diagnostic = problem.get('error', {}).get('message', '') or str(problem.get('issue', []))
            assert response.status == status, f'{method} {path}: expected {status}, got {response.status}; {diagnostic}'
            if fhir:
                assert response.headers.get_content_type() == 'application/fhir+json', f'{path}: wrong FHIR media type'
        return json.loads(data) if data else None

    setup = call('GET', '/api/v1/auth/setup')
    assert not setup['configured'], 'Use a fresh disposable stack, never the main HIS database'
    password = 'Aa1' + secrets.token_urlsafe(24)
    admin = call('POST', '/api/v1/auth/setup', {'password': password, 'confirmation': password, 'secret': setup['secret'], 'totp': totp(setup['secret'])})['access_token']
    print('OK: isolated administrator setup; no seeded administrator')
    workers = {}
    for role in ['doctor', 'nurse', 'laboratory-technician', 'other']:
        staff = call('POST', '/api/v1/staff', {'first_name': 'Тест', 'last_name': 'Лаборатория', 'position': 'Врач', 'role': role}, token=admin, status=201)
        account = call('POST', '/api/v1/accounts', {'staff_id': staff['id'], 'username': 'u' + uuid.uuid4().hex[:10]}, token=admin, status=201)
        session = call('POST', '/api/v1/auth/login', {'username': account['account']['username'], 'password': account['password']})
        workers[role] = (staff, session['access_token'])
    doctor, token = workers['doctor']
    patient = call('POST', '/api/v1/patients', {'first_name': 'Учебный', 'last_name': 'Пациент', 'birth_date': '1990-01-01', 'administrative_sex': 'M'}, token=token, status=201)
    call('POST', f'/api/v1/patients/{patient["id"]}/record', {'staff_id': doctor['id']}, token=token, status=201)
    tests = call('GET', '/api/v1/lab/tests', token=token)
    assert {t['code'] for t in tests} == {'2345-7', '33706-3', '718-7', '56888-1', '5195-3', '13955-0', '22587-0'}
    extra = call('POST', '/api/v1/lab/tests', {'code': '4548-4', 'name': 'Гликированный гемоглобин'}, token=admin, status=201)
    call('PUT', f'/api/v1/lab/tests/{extra["id"]}', {'code': '4548-4', 'name': 'Гликированный гемоглобин (HbA1c)'}, token=admin)
    call('DELETE', f'/api/v1/lab/tests/{extra["id"]}', token=admin)
    call('POST', '/api/v1/lab/orders', {'patient_id': patient['id'], 'test_ids': [extra['id']]}, token=token, status=400)
    for role, (_, worker_token) in workers.items():
        if role != 'doctor':
            call('POST', '/api/v1/lab/orders', {'patient_id': patient['id'], 'test_ids': [tests[0]['id']]}, token=worker_token, status=403)
    orders = call('POST', '/api/v1/lab/orders', {'patient_id': patient['id'], 'test_ids': [t['id'] for t in tests], 'comment': 'Учебная проверка FHIR R5'}, token=token, status=201)
    requests = [e['resource'] for e in orders['entry']]
    assert len(requests) == 7 and len({r['id'] for r in requests}) == 7
    assert all(len(r['code']['concept']['coding']) == 1 for r in requests)
    repeated = call('POST', '/api/v1/lab/orders', {'patient_id': patient['id'], 'test_ids': [tests[0]['id'], tests[0]['id']]}, token=token, status=201)
    assert repeated['entry'][0]['resource']['id'] != repeated['entry'][1]['resource']['id']
    print('OK: LOINC directory, multiple independent orders, doctor-only ordering')
    assert not call('GET', '/lis/v1/connection', lis=True)['configured']
    call('GET', '/lis/v1/requests', lis=True, status=409)
    service = call('POST', '/api/v1/service-accounts', {'name': 'lab-integration'}, token=admin, status=201)
    key = service['api_key']
    service_id = service['account']['id']
    assert service['account']['role'] == 'lis'
    assert all(a['id'] != service_id for a in call('GET', '/api/v1/accounts', token=admin)['items'])
    call('POST', '/api/v1/auth/login', {'username': 'lab-integration', 'password': key}, status=401)
    call('POST', '/api/v1/auth/service', {'api_key': 'x' * 64}, status=401)
    session = call('POST', '/api/v1/auth/service', {'api_key': key})
    previous_refresh = session['refresh_token']
    session = call('POST', '/api/v1/auth/refresh', {'refresh_token': previous_refresh})
    call('POST', '/api/v1/auth/refresh', {'refresh_token': previous_refresh}, status=401)
    access = session['access_token']
    call('PUT', '/lis/v1/connection', {'api_key': 'x' * 64}, lis=True, status=502)
    assert not call('GET', '/lis/v1/connection', lis=True)['configured']
    call('PUT', '/lis/v1/connection', {'address': args.his, 'api_key': key}, lis=True, status=400)
    status = call('PUT', '/lis/v1/connection', {'api_key': key}, lis=True)
    assert status['configured'] and 'api_key' not in status and 'address' not in status
    assert call('GET', '/lis/v1/connection', lis=True)['configured']
    for path in ['/api/v1/accounts', '/api/v1/staff', '/api/v1/patients', '/api/v1/lab/tests', '/api/v1/service-accounts']:
        call('GET', path, token=access, status=403)
    for path in ['/api/v1/fhir/ServiceRequest', '/api/v1/fhir/Observation']:
        bundle = call('GET', path, token=access, fhir=True)
        assert bundle['resourceType'] == 'Bundle'
    call('POST', '/api/v1/fhir/ServiceRequest', {}, token=access, status=403, fhir=True)
    print('OK: API-key login, tokens, separate accounts, restricted LIS role, TLS verification')
    queue = call('GET', '/lis/v1/requests', lis=True)
    assert queue['total'] == 9
    first = requests[0]
    details = call('GET', f'/lis/v1/requests/{first["id"]}', lis=True)
    assert details['patient']['resourceType'] == 'Patient'
    assert details['patient']['id'] == patient['id']
    assert details['practitioner']['resourceType'] == 'Practitioner'
    assert details['practitioner']['id'] == doctor['id']
    assert not call('GET', '/lis/v1/history', lis=True), 'LIS must not persist unprocessed requests'
    now = datetime.now(timezone.utc).isoformat()
    def observation(order, value='142'):
        return {'resourceType': 'Observation', 'id': str(uuid.uuid4()), 'status': 'final', 'code': order['code']['concept'], 'subject': order['subject'], 'basedOn': [{'reference': 'ServiceRequest/' + order['id']}], 'effectiveDateTime': now, 'issued': now, 'valueString': value}
    invalid = observation(first)
    invalid['subject'] = {'reference': 'Patient/' + str(uuid.uuid4())}
    problem = call('POST', '/api/v1/fhir/Observation', invalid, token=access, status=400, fhir=True)
    assert problem['resourceType'] == 'OperationOutcome'
    for _, worker_token in workers.values():
        call('POST', '/api/v1/fhir/Observation', observation(first), token=worker_token, status=403, fhir=True)
    for index, order in enumerate(requests):
        history = call('POST', '/lis/v1/results', {'request_id': order['id'], 'value': f'Значение {index + 1}'}, lis=True, status=201)
        assert history['delivered'] and history['result']['valueString'] == f'Значение {index + 1}'
        result = history['result']
        assert result['basedOn'] == [{'reference': 'ServiceRequest/' + order['id']}]
        assert result['subject']['reference'] == 'Patient/' + patient['id']
        call('POST', '/lis/v1/results', {'request_id': order['id'], 'value': 'Повтор'}, lis=True, status=409)
        details = call('GET', f'/lis/v1/requests/{order["id"]}', lis=True)
        assert details['history']['result']['id'] == result['id']
        assert call('GET', '/api/v1/fhir/Observation/' + result['id'], token=access, fhir=True)['valueString'] == result['valueString']
        call('PUT', '/api/v1/fhir/Observation/' + result['id'], result, token=access, status=403, fhir=True)
    for _, worker_token in workers.values():
        assert call('GET', '/api/v1/fhir/Observation?patient=' + patient['id'], token=worker_token, fhir=True)['total'] == 7
    assert len(call('GET', '/lis/v1/history?search=%D0%9F%D0%B0%D1%86%D0%B8%D0%B5%D0%BD%D1%82', lis=True)) == 7
    assert call('GET', '/api/v1/fhir/ServiceRequest/' + first['id'], token=access, fhir=True)['status'] == 'completed'
    call('POST', '/api/v1/fhir/Observation', observation(first), token=access, status=409, fhir=True)
    print('OK: seven results, FHIR references, HIS visibility, immutable LIS snapshots, duplicate protection')
    parallel_order = repeated['entry'][0]['resource']
    def submit(_):
        payload = observation(parallel_order)
        request = urllib.request.Request(args.his + '/api/v1/fhir/Observation', data=json.dumps(payload).encode(), headers={'Content-Type': 'application/fhir+json', 'Authorization': 'Bearer ' + access}, method='POST')
        try:
            with urllib.request.urlopen(request, context=context, timeout=25) as response:
                return response.status
        except urllib.error.HTTPError as error:
            return error.code
    with ThreadPoolExecutor(max_workers=4) as executor:
        statuses = list(executor.map(submit, range(4)))
    assert sorted(statuses) == [201, 409, 409, 409]
    foreign = call('GET', '/lis/v1/requests/' + parallel_order['id'], lis=True)
    assert foreign['result']['basedOn'][0]['reference'] == 'ServiceRequest/' + parallel_order['id']
    assert 'history' not in foreign
    assert len(call('GET', '/lis/v1/history', lis=True)) == 7
    rotated = call('POST', f'/api/v1/service-accounts/{service_id}/key', token=admin)
    call('POST', '/api/v1/auth/refresh', {'refresh_token': session['refresh_token']}, status=401)
    call('POST', '/api/v1/auth/service', {'api_key': key}, status=401)
    call('GET', '/api/v1/fhir/ServiceRequest', token=access, status=401, fhir=True)
    call('GET', '/lis/v1/requests', lis=True, status=502)
    assert len(call('GET', '/lis/v1/history', lis=True)) == 7
    new_session = call('POST', '/api/v1/auth/service', {'api_key': rotated['api_key']})
    call('PUT', f'/api/v1/service-accounts/{service_id}', {'name': 'lab-integration', 'is_active': False}, token=admin)
    call('POST', '/api/v1/auth/refresh', {'refresh_token': new_session['refresh_token']}, status=401)
    call('GET', '/api/v1/fhir/ServiceRequest', token=new_session['access_token'], status=401, fhir=True)
    call('POST', '/api/v1/auth/service', {'api_key': rotated['api_key']}, status=401)
    print('OK: concurrent uniqueness, key rotation, deactivation, revoked access/refresh, offline local history')
    print('HIS/LIS end-to-end checks passed.')


if __name__ == '__main__':
    main()
