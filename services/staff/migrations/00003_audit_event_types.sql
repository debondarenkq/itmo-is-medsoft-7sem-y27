-- +goose Up
ALTER TABLE audit_events ADD COLUMN event_type text NOT NULL DEFAULT 'other';
UPDATE audit_events SET event_type=CASE
 WHEN message LIKE 'Настройка администратора%' THEN 'admin_setup'
 WHEN message LIKE '%Неуспеш%' OR message LIKE '%Превышен лимит%' THEN 'login_failed'
 WHEN message LIKE 'Вход %' THEN 'login'
 WHEN message LIKE 'Выход %' THEN 'logout'
 WHEN message LIKE 'Обновление пары%' THEN 'token_refreshed'
 WHEN message LIKE 'Отказ в доступе%' OR message LIKE 'Неавторизованный%' THEN 'access_denied'
 WHEN message LIKE 'Создание аккаунта%' THEN 'account_created'
 WHEN message LIKE 'Изменение аккаунта%' OR message LIKE 'Изменение роли/%' THEN 'account_updated'
 WHEN message LIKE '%POST /api/v1/patients/%/record%' THEN 'record_created'
 WHEN message LIKE '%POST /api/v1/patients%' THEN 'patient_created'
 WHEN message LIKE '%PUT /api/v1/patients/%' THEN 'patient_updated'
 WHEN message LIKE '%DELETE /api/v1/patients/%' THEN 'patient_deleted'
 WHEN message LIKE '%POST /api/v1/records/%' THEN 'record_updated'
 WHEN message LIKE '%POST /api/v1/staff%' THEN 'staff_created'
 WHEN message LIKE '%PUT /api/v1/staff/%' THEN 'staff_updated'
 WHEN message LIKE '%DELETE /api/v1/staff/%' THEN 'staff_deleted'
 ELSE 'other' END;
CREATE INDEX audit_event_type_time ON audit_events(event_type,event_time,event_id);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION staff_security_change() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE affected uuid; revoked integer;
BEGIN
 IF NEW.role IS DISTINCT FROM OLD.role OR (NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL) THEN
  SELECT id INTO affected FROM accounts WHERE staff_id=NEW.id FOR UPDATE;
  IF affected IS NOT NULL THEN
   IF NEW.deleted_at IS NOT NULL THEN
    UPDATE accounts SET is_active=false, updated_at=clock_timestamp() WHERE id=affected;
   END IF;
   DELETE FROM refresh_tokens WHERE account_id=affected;
   GET DIAGNOSTICS revoked = ROW_COUNT;
   INSERT INTO audit_events(event_id,actor,message,event_type) VALUES(gen_random_uuid(), 'system:staff:'||NEW.id,
    'Изменение роли/удаление медработника; аккаунт '||affected||'; отозвано refresh токенов: '||revoked,'account_updated');
  END IF;
 END IF;
 RETURN NEW;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION staff_security_change() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE affected uuid; revoked integer;
BEGIN
 IF NEW.role IS DISTINCT FROM OLD.role OR (NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL) THEN
  SELECT id INTO affected FROM accounts WHERE staff_id=NEW.id FOR UPDATE;
  IF affected IS NOT NULL THEN
   IF NEW.deleted_at IS NOT NULL THEN
    UPDATE accounts SET is_active=false, updated_at=clock_timestamp() WHERE id=affected;
   END IF;
   DELETE FROM refresh_tokens WHERE account_id=affected;
   GET DIAGNOSTICS revoked = ROW_COUNT;
   INSERT INTO audit_events(event_id,actor,message) VALUES(gen_random_uuid(), 'system:staff:'||NEW.id,
    'Изменение роли/удаление медработника; аккаунт '||affected||'; отозвано refresh токенов: '||revoked);
  END IF;
 END IF;
 RETURN NEW;
END $$;
-- +goose StatementEnd
ALTER TABLE audit_events DROP COLUMN event_type;
