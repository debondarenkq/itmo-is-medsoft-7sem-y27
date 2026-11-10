#!/usr/bin/env python3
import argparse
from datetime import datetime, timedelta, timezone
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import uuid

ROOT = Path(__file__).resolve().parent.parent
COMPOSE = ROOT / 'scripts/compose.sh'
DATABASES = ('identity', 'staff', 'catalog', 'clinical')
APIS = ('identity-service', 'staff-service', 'catalog-service', 'clinical-service')


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


def reset_database(database):
    query(database, f'DROP SCHEMA public CASCADE; CREATE SCHEMA public AUTHORIZATION {database}_user;')


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


def private_file(path, content):
    path.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.NamedTemporaryFile(dir=path.parent, delete=False) as stream:
        temporary = Path(stream.name)
        try:
            os.chmod(temporary, 0o600)
            stream.write(content)
            stream.flush()
            os.replace(temporary, path)
        finally:
            temporary.unlink(missing_ok=True)


def main():
    parser = argparse.ArgumentParser(description='Reset local data and load demo fixtures without accounts.')
    parser.parse_args()
    staff, diagnoses, patients, records, events, audit = demo_data()
    subprocess.run(['make', 'certs'], cwd=ROOT, check=True)
    compose('up', '--no-build', '--pull', 'never', '-d', '--wait', 'postgres')
    configuration = json.loads(compose('config', '--format', 'json', text=True, capture_output=True).stdout)
    mounts = configuration['services']['identity-service'].get('volumes', [])
    keys = [mount['source'] for mount in mounts
            if mount.get('type') == 'bind' and mount.get('target') == '/security/jwt/signing.key']
    if len(keys) != 1:
        raise RuntimeError('Нужен файловый bind mount ключа JWT в identity-service.')
    signing_key = Path(keys[0])
    project = configuration['name']
    output = ROOT / 'build' if project == 'his' else ROOT / 'build' / project
    print('Reset: очищаем данные четырёх сервисов и загружаем демонстрационный набор.', flush=True)
    compose('stop', 'nginx', *APIS)
    try:
        for database in DATABASES:
            reset_database(database)
        key = subprocess.run(['openssl', 'genpkey', '-algorithm', 'ED25519'],
                             check=True, capture_output=True).stdout
        private_file(signing_key, key)
        compose('up', '--no-build', '--pull', 'never', '--no-deps', '--force-recreate', '-d', '--wait', *APIS)
        query('staff', 'BEGIN;\n' + insert('staff', staff) + '\nCOMMIT;')
        query('identity', 'BEGIN;\n' + insert('audit_events', audit) + '\nCOMMIT;')
        query('catalog', 'BEGIN;\n' + insert('diagnoses', diagnoses) + '\nCOMMIT;')
        query('clinical', 'BEGIN;\n' + insert('patients', patients) + insert('records', records) + insert('record_events', events) + '\nCOMMIT;')
        (output / 'demo-accounts.json').unlink(missing_ok=True)
        if query('identity', 'SELECT count(*) FROM accounts WHERE is_admin') != '0':
            raise RuntimeError('Reset не должен создавать администратора.')
        compose('up', '--no-build', '--pull', 'never', '-d', '--wait', 'nginx')
    except Exception:
        print('Reset не завершён; gateway оставлен остановленным. Повторите make reset после устранения ошибки.', file=sys.stderr)
        raise
    print('Готово: 5 сотрудников (4 активных), 6 диагнозов, 6 пациентов, 3 ЭМК с историей.')
    print('Администратор не создан. Настройте пароль и TOTP при первом входе.')
    print('Демо-аккаунты: после настройки администратора запустите python3 scripts/create-demo-accounts.py.')
    print('Предыдущие сеансы недействительны; в открытом клиенте выполните новый вход.')


if __name__ == '__main__':
    try:
        main()
    except subprocess.CalledProcessError as error:
        if error.stderr:
            print(error.stderr.decode() if isinstance(error.stderr, bytes) else error.stderr, file=sys.stderr)
        sys.exit(f'Команда завершилась ошибкой: {error.cmd[0]} (код {error.returncode})')
    except Exception as error:
        sys.exit(str(error))
