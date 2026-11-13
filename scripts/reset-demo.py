#!/usr/bin/env python3
import argparse
from datetime import datetime, timedelta, timezone
import json
from pathlib import Path
import subprocess
import sys
import uuid

ROOT = Path(__file__).resolve().parent.parent
COMPOSE = ROOT / 'scripts/compose.sh'
APIS = ('identity-service', 'staff-service', 'catalog-service', 'clinical-service', 'laboratory-service')


def compose(*args, **kwargs):
    return subprocess.run([str(COMPOSE), *args], cwd=ROOT, check=True, **kwargs)


def query(database, sql):
    return compose('exec', '-T', 'postgres', 'psql', '-X', '-qAt', '-v',
                   'ON_ERROR_STOP=1', '-U', 'his_admin', '-d', database,
                   input=sql, text=True, capture_output=True).stdout.strip()


def value(item):
    if item is None:
        return 'NULL'
    if isinstance(item, bool):
        return 'TRUE' if item else 'FALSE'
    if isinstance(item, (int, float)):
        return str(item)
    if isinstance(item, (dict, list)):
        item = json.dumps(item, ensure_ascii=False)
    return "'" + str(item).replace("'", "''") + "'"


def insert(table, rows):
    if not rows:
        return ''
    columns = list(rows[0])
    items = ',\n'.join('(' + ','.join(value('null' if row[key] is None and key in ('before_data', 'after_data') else row[key]) for key in columns) + ')' for row in rows)
    return f"INSERT INTO {table} ({','.join(columns)}) VALUES {items};\n"


def demo_data():
    now = datetime.now(timezone.utc)
    def at(days):
        return (now - timedelta(days=days)).isoformat()
    def identifier(label):
        return str(uuid.uuid5(uuid.NAMESPACE_URL, 'his-demo/' + label))

    staff = []
    profiles = json.loads((ROOT / 'examples/demo-staff.json').read_text())
    for index, profile in enumerate(profiles):
        staff.append(dict(id=identifier(f'staff/{index}'),
                          **{key: profile[key] for key in ('first_name', 'last_name', 'position', 'role')},
                          created_at=at(14), updated_at=at(1),
                          deleted_at=at(1) if profile['deleted'] else None))
    diagnoses = []
    for index, (code, name, deleted) in enumerate([
        ('I10.00', 'Гипертензия', False), ('J06.90', 'ОРВИ', False),
        ('E11.90', 'Сахарный диабет', False), ('J45.90', 'Бронхиальная астма', False),
        ('K29.50', 'Хронический гастрит', False), ('M54.50', 'Боль в пояснице', True),
    ]):
        diagnoses.append(dict(id=identifier(f'diagnosis/{index}'), code=code, name=name,
                              created_at=at(14), updated_at=at(1), deleted_at=at(1) if deleted else None))
    patients = []
    for index, (first, last, middle, birth, sex) in enumerate([
        ('Сергей', 'Иванов', 'Петрович', '1985-04-12', 'M'),
        ('Анна', 'Петрова', 'Сергеевна', '1989-08-23', 'F'),
        ('Алексей', 'Смирнов', None, '2001-02-17', 'M'),
        ('Мария', 'Кузнецова', 'Ивановна', '1994-11-05', 'F'),
        ('Дмитрий', 'Волков', 'Андреевич', '1970-06-30', 'M'),
        ('Саша', 'Новикова', None, '2005-01-10', 'UNKNOWN'),
    ]):
        patients.append(dict(id=identifier(f'patient/{index}'), first_name=first, last_name=last,
                             middle_name=middle, birth_date=birth, administrative_sex=sex,
                             comment='Учебные данные', created_at=at(7), updated_at=at(1)))
    actor = {key: staff[0][key] for key in ('id', 'first_name', 'last_name', 'position')}
    records, events = [], []
    for index, patient in enumerate(patients[:3]):
        rid = identifier(f'record/{index}')
        state = dict(diagnoses=[], prescriptions=[])
        def emit(kind, entity, before, after, version, days):
            events.append(dict(id=identifier(f'event/{index}/{len(events)}'), record_id=rid,
                               sequence=1 + sum(event['record_id'] == rid for event in events),
                               version=version, occurred_at=at(days), actor=actor, event_type=kind,
                               entity_id=entity, before_data=before, after_data=after))
        emit('record_created', rid, None, dict(id=rid, patient_id=patient['id']), 1, 5)
        if index < 2:
            diagnosis = diagnoses[index]
            entry = dict(id=identifier(f'record-diagnosis/{index}'), diagnosis_id=diagnosis['id'],
                         code=diagnosis['code'], name=diagnosis['name'])
            prescription = dict(id=identifier(f'prescription/{index}'),
                                text='Контроль артериального давления' if index == 0 else 'Контроль температуры',
                                status='active')
            state['diagnoses'].append(entry)
            emit('diagnosis_added', entry['id'], None, entry, 2, 4)
            state['prescriptions'].append(prescription)
            emit('prescription_added', prescription['id'], None, prescription, 2, 4)
            updated = {**prescription, **(dict(text='Контроль давления утром и вечером')
                                         if index == 0 else dict(status='cancelled'))}
            emit('prescription_updated' if index == 0 else 'prescription_cancelled',
                 prescription['id'], prescription, updated, 3, 2)
            state['prescriptions'][0] = updated
        records.append(dict(id=rid, patient_id=patient['id'], version=3 if index < 2 else 1,
                            state=state, created_at=at(5), updated_at=at(2 if index < 2 else 5)))
    audit = [dict(event_id=identifier('audit/' + kind), event_time=at(0), actor='system:demo',
                  message=message, event_type=kind) for kind, message in [
        ('staff_created', 'Подготовлены демонстрационные медработники'),
        ('patient_created', 'Подготовлены демонстрационные пациенты'),
        ('record_created', 'Подготовлены демонстрационные ЭМК и история'),
    ]]
    return staff, diagnoses, patients, records, events, audit


def seed_lab_requests(staff, patients, records):
    tests = json.loads(query('laboratory', "SELECT json_agg(t ORDER BY t.name,t.id) FROM lab_tests t WHERE active"))
    doctor = next(person for person in staff if person['role'] == 'doctor' and person['deleted_at'] is None)
    by_id = {patient['id']: patient for patient in patients}
    now = datetime.now(timezone.utc).isoformat()
    requests = []
    for index, test in enumerate(tests):
        patient = by_id[records[index % len(records)]['patient_id']]
        patient_name = ' '.join(part for part in (patient['last_name'], patient['first_name'], patient['middle_name']) if part)
        request_id = str(uuid.uuid5(uuid.NAMESPACE_URL, f'his-demo/lab-request/{index}'))
        resource = {
            'resourceType': 'ServiceRequest', 'id': request_id,
            'status': 'active', 'intent': 'order', 'authoredOn': now,
            'code': {'concept': {'coding': [{'system': 'http://loinc.org', 'code': test['code'], 'display': test['name']}]}},
            'subject': {'reference': 'Patient/' + patient['id'], 'display': patient_name},
            'requester': {'reference': 'Practitioner/' + doctor['id'], 'display': doctor['last_name'] + ' ' + doctor['first_name']},
            'note': [{'text': 'Демонстрационное назначение'}],
        }
        requests.append(dict(id=request_id, patient_id=patient['id'], created_at=now, resource=resource))
    query('laboratory', 'BEGIN;\n' + insert('lab_requests', requests) + '\nCOMMIT;')
    return len(requests)


def main():
    parser = argparse.ArgumentParser(description='Recreate HIS and LIS databases; optionally load demo data without accounts.')
    parser.add_argument('--empty', action='store_true', help='Start with empty business data and initial administrator setup.')
    args = parser.parse_args()
    subprocess.run(['make', 'certs', 'gateway-network'], cwd=ROOT, check=True)
    configuration = json.loads(compose('config', '--format', 'json', text=True, capture_output=True).stdout)
    project = configuration['name']
    output = ROOT / 'build' if project == 'his' else ROOT / 'build' / project
    print('Полностью пересоздаём данные HIS и LIS.', flush=True)
    compose('-f', 'deploy/lis.compose.yaml', 'down', '-v', '--remove-orphans')
    compose('down', '-v', '--remove-orphans')
    compose('up', '--no-build', '--pull', 'never', '-d', '--wait', *APIS)
    if not args.empty:
        staff, diagnoses, patients, records, events, audit = demo_data()
        query('staff', 'BEGIN;\n' + insert('staff', staff) + '\nCOMMIT;')
        query('identity', 'BEGIN;\n' + insert('audit_events', audit) + '\nCOMMIT;')
        query('catalog', 'BEGIN;\n' + insert('diagnoses', diagnoses) + '\nCOMMIT;')
        query('clinical', 'BEGIN;\n' + insert('patients', patients) + insert('records', records) + insert('record_events', events) + '\nCOMMIT;')
        count = seed_lab_requests(staff, patients, records)
        print(f'Загружены 5 сотрудников, 6 диагнозов, 6 пациентов, 3 ЭМК и {count} назначений анализов.')
    (output / 'demo-accounts.json').unlink(missing_ok=True)
    if query('identity', 'SELECT count(*) FROM accounts') != '0':
        raise RuntimeError('Начальное заполнение не должно создавать аккаунты.')
    compose('up', '--no-build', '--pull', 'never', '-d', '--wait', 'nginx')
    compose('-f', 'deploy/lis.compose.yaml', 'up', '--no-build', '--pull', 'never', '-d', '--wait')
    print('HIS и LIS запущены. Настройте пароль и TOTP администратора в HISAdmin.')
    if not args.empty:
        print('После настройки администратора make demo-accounts создаст демо-аккаунты медработников через API.')
    print('Для LIS создайте сервисный аккаунт и введите его API-ключ.')


if __name__ == '__main__':
    try:
        main()
    except subprocess.CalledProcessError as error:
        if error.stderr:
            print(error.stderr.decode() if isinstance(error.stderr, bytes) else error.stderr, file=sys.stderr)
        sys.exit(f'Команда завершилась ошибкой: {error.cmd[0]} (код {error.returncode})')
    except Exception as error:
        sys.exit(str(error))
