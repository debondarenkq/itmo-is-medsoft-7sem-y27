#!/usr/bin/env python3
"""Exercise all HIS scenarios through the public gateway; uses synthetic data.

Creates demo records intentionally: records cannot be deleted through the API.
Run against a disposable Compose project when checking a clean database.
"""
import argparse
import datetime as dt
import json
from pathlib import Path
import urllib.error
import urllib.parse
import urllib.request
import uuid


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--url', default='http://127.0.0.1:8080')
    args = parser.parse_args()
    base = args.url.rstrip('/')

    def call(method, path, body=None, status=200):
        payload = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request(base + path, data=payload, method=method)
        if payload is not None:
            request.add_header('Content-Type', 'application/json')
        try:
            response = urllib.request.urlopen(request, timeout=20)
        except urllib.error.HTTPError as error:
            response = error
        with response:
            actual = response.status
            data = response.read()
        assert actual == status, f'{method} {path}: expected {status}, got {actual}: {data.decode()}'
        return json.loads(data) if data else None

    for service in ('staff', 'catalog', 'clinical'):
        spec = call('GET', f'/docs/{service}/openapi.json')
        assert spec['openapi'] == '3.0.3'

    file_data = json.loads((Path(__file__).resolve().parent.parent / 'examples' / 'diagnoses.json').read_text())
    imported = call('POST', '/api/v1/diagnoses/import', file_data)['items']
    imported_again = call('POST', '/api/v1/diagnoses/import', file_data)['items']
    assert len(imported) == 5 and [d['id'] for d in imported] == [d['id'] for d in imported_again]
    call('POST', '/api/v1/diagnoses/import', {'diagnoses': []}, 400)
    call('POST', '/api/v1/diagnoses/import', {'diagnoses': [{'code': 'E11.900', 'name': 'Диагноз'}]}, 400)
    call('POST', '/api/v1/diagnoses/import', {'diagnoses': [file_data['diagnoses'][0], file_data['diagnoses'][0]]}, 400)
    call('POST', '/api/v1/staff', {'first_name': 'А', 'last_name': 'Иванов', 'position': 'Врач'}, 400)
    staff_input = {'first_name': 'Анна', 'last_name': 'Петрова', 'position': 'Терапевт'}
    staff = call('POST', '/api/v1/staff', staff_input, 201)
    sid = staff['id']
    assert call('GET', f'/api/v1/staff/{sid}')['first_name'] == 'Анна'
    staff_input['position'] = 'Доктор'
    staff = call('PUT', f'/api/v1/staff/{sid}', staff_input)

    code = uuid.uuid4().hex[:6].upper()
    diagnosis_input = {'code': code, 'name': 'Диагноз ' + uuid.uuid4().hex[:8]}
    diagnosis = call('POST', '/api/v1/diagnoses', diagnosis_input, 201)
    did = diagnosis['id']
    assert call('GET', f'/api/v1/diagnoses/{did}')['code'] == code
    call('POST', '/api/v1/diagnoses', {'code': code, 'name': 'Другая запись'}, 409)
    call('POST', '/api/v1/diagnoses', {'code': uuid.uuid4().hex[:6].upper(), 'name': diagnosis['name']}, 409)
    call('POST', '/api/v1/diagnoses', {'code': 'short', 'name': 'Название'}, 400)

    patient_input = {'first_name': 'Иван', 'last_name': 'Иванов', 'middle_name': None,
                     'birth_date': '2000-02-29', 'administrative_sex': 'UNKNOWN', 'comment': ''}
    call('POST', '/api/v1/patients', {**patient_input, 'birth_date': '3000-01-01'}, 400)
    call('POST', '/api/v1/patients', {**patient_input, 'comment': 'я' * 33}, 400)
    disposable = call('POST', '/api/v1/patients', patient_input, 201)
    call('DELETE', f'/api/v1/patients/{disposable["id"]}', status=204)
    call('GET', f'/api/v1/patients/{disposable["id"]}', status=404)
    patient = call('POST', '/api/v1/patients', patient_input, 201)
    pid = patient['id']
    patient_input['comment'] = 'Учебный пациент'
    patient = call('PUT', f'/api/v1/patients/{pid}', patient_input)
    assert not patient['has_record']
    call('POST', f'/api/v1/patients/{pid}/record', {}, 400)
    record = call('POST', f'/api/v1/patients/{pid}/record', {'staff_id': sid}, 201)
    rid = record['id']
    path = f'/api/v1/records/{rid}'
    assert record['version'] == 1 and record['state'] == {'diagnoses': [], 'prescriptions': []}
    call('POST', f'/api/v1/patients/{pid}/record', {'staff_id': sid}, 409)
    assert call('GET', f'/api/v1/patients/{pid}/record')['id'] == rid
    call('DELETE', f'/api/v1/patients/{pid}', status=409)
    assert call('GET', f'/api/v1/patients/{pid}')['has_record']
    call('POST', path + '/changes', {'expected_version': 1, 'commands': [{'type': 'add_prescription', 'text': 'Пить воду'}]}, 400)

    def save(version, commands, status=200, author=sid):
        return call('POST', path + '/changes', {'staff_id': author, 'expected_version': version, 'commands': commands}, status)

    save(1, [{'type': 'add_prescription', 'text': 'Пить воду'}, {'type': 'remove_diagnosis', 'id': str(uuid.uuid4())}], 404)
    assert call('GET', path)['version'] == 1
    assert call('GET', path + '/history')['total'] == 1
    added = save(1, [{'type': 'add_diagnosis', 'diagnosis_id': did}, {'type': 'add_prescription', 'text': 'Пить больше воды'}])
    assert added['version'] == 2
    old_rx = added['state']['prescriptions'][0]['id']
    entry_id = added['state']['diagnoses'][0]['id']
    save(2, [{'type': 'add_diagnosis', 'diagnosis_id': did}], 409)
    save(1, [{'type': 'add_prescription', 'text': 'Назначение'}], 409)
    staff_input['first_name'] = 'Мария'
    call('PUT', f'/api/v1/staff/{sid}', staff_input)
    diagnosis_input['name'] = 'Изменённый ' + uuid.uuid4().hex[:6]
    call('PUT', f'/api/v1/diagnoses/{did}', diagnosis_input)
    edited = save(2, [{'type': 'edit_prescription', 'id': old_rx, 'text': 'Пить воду после еды'}])
    new_rx = edited['state']['prescriptions'][0]['id']
    assert new_rx != old_rx and edited['version'] == 3
    assert edited['state']['diagnoses'][0]['name'] == diagnosis['name']
    cancelled = save(3, [{'type': 'cancel_prescription', 'id': new_rx}, {'type': 'remove_diagnosis', 'id': entry_id}])
    assert cancelled['version'] == 4
    assert cancelled['state']['diagnoses'] == []
    assert cancelled['state']['prescriptions'][0]['status'] == 'cancelled'
    save(4, [{'type': 'cancel_prescription', 'id': new_rx}], 409)
    save(4, [{'type': 'edit_prescription', 'id': old_rx, 'text': 'Новое назначение'}], 404)
    history = call('GET', path + '/history')['items']
    assert len(history) == 7
    assert [e['sequence'] for e in history] == list(range(1, 8))
    assert history[1]['actor']['first_name'] == 'Анна'
    assert history[3]['actor']['first_name'] == 'Мария'
    assert [e['type'] for e in history[3:5]] == ['prescription_removed', 'prescription_added']
    assert history[3]['occurred_at'] == history[4]['occurred_at']
    assert history[3]['version'] == history[4]['version'] == 3
    assert history[3]['before']['text'] == 'Пить больше воды'
    for snapshot in (added, edited, cancelled):
        restored = call('GET', path + '/state?' + urllib.parse.urlencode({'at': snapshot['updated_at']}))
        assert restored['state'] == snapshot['state'] and restored['version'] == snapshot['version']
        restored = call('GET', path + f'/state?version={snapshot["version"]}')
        assert restored['state'] == snapshot['state']
    just_before_edit = dt.datetime.fromisoformat(edited['updated_at'].replace('Z', '+00:00')) - dt.timedelta(microseconds=1)
    restored = call('GET', path + '/state?' + urllib.parse.urlencode({'at': just_before_edit.isoformat()}))
    assert restored['state'] == added['state']
    call('GET', path + '/state?at=1900-01-01T00:00:00Z', status=404)
    call('GET', path + '/state?version=999', status=404)
    call('GET', path + '/state?version=1&at=2000-01-01T00:00:00Z', status=400)
    page = call('GET', path + '/history?limit=2&offset=3')
    assert len(page['items']) == 2 and page['total'] == 7

    call('DELETE', f'/api/v1/staff/{sid}', status=204)
    assert call('GET', f'/api/v1/staff/{sid}')['deleted_at'] is not None
    assert all(s['id'] != sid for s in call('GET', '/api/v1/staff?' + urllib.parse.urlencode({'q': 'Петрова'}))['items'])
    save(4, [{'type': 'add_prescription', 'text': 'Назначение'}], 400)
    second_staff = call('POST', '/api/v1/staff', staff_input, 201)
    call('DELETE', f'/api/v1/diagnoses/{did}', status=204)
    save(4, [{'type': 'add_diagnosis', 'diagnosis_id': did}], 400, second_staff['id'])
    assert call('GET', path)['version'] == 4
    call('DELETE', f'/api/v1/staff/{second_staff["id"]}', status=204)
    print('PASS: staff, diagnoses, patients, EMR editing, history, snapshots, replay, rollback and API constraints')


if __name__ == '__main__':
    main()
