-- +goose Up
ALTER TABLE staff DROP CONSTRAINT staff_first_name_check;
ALTER TABLE staff DROP CONSTRAINT staff_last_name_check;
ALTER TABLE staff ADD CHECK (char_length(first_name) BETWEEN 1 AND 24);
ALTER TABLE staff ADD CHECK (char_length(last_name) BETWEEN 1 AND 24);
ALTER TABLE staff ADD COLUMN role text NOT NULL DEFAULT 'other'
 CHECK (role IN ('doctor','laboratory-technician','nurse','other'));
CREATE TABLE accounts (
 id uuid PRIMARY KEY,
 username text NOT NULL UNIQUE CHECK (char_length(username) BETWEEN 3 AND 16),
 password_hash text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 is_active boolean NOT NULL DEFAULT true,
 staff_id uuid UNIQUE REFERENCES staff(id),
 is_admin boolean NOT NULL DEFAULT false,
 totp_secret text,
 last_totp_step bigint NOT NULL DEFAULT -1,
 CHECK ((is_admin AND staff_id IS NULL AND username='admin') OR (NOT is_admin))
);
CREATE UNIQUE INDEX one_admin ON accounts (is_admin) WHERE is_admin;
CREATE TABLE refresh_tokens (
 token_hash text PRIMARY KEY,
 account_id uuid NOT NULL REFERENCES accounts(id),
 expires_at timestamptz NOT NULL
);
CREATE INDEX refresh_account ON refresh_tokens(account_id);
CREATE TABLE audit_events (
 event_id uuid PRIMARY KEY,
 event_time timestamptz NOT NULL DEFAULT clock_timestamp(),
 actor text NOT NULL,
 message text NOT NULL
);
CREATE INDEX audit_time ON audit_events(event_time, event_id);
CREATE TABLE security_settings (id boolean PRIMARY KEY DEFAULT true CHECK(id), setup_secret text NOT NULL);
-- +goose StatementBegin
CREATE FUNCTION staff_security_change() RETURNS trigger LANGUAGE plpgsql AS $$
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
CREATE TRIGGER staff_security AFTER UPDATE ON staff FOR EACH ROW EXECUTE FUNCTION staff_security_change();

-- +goose Down
DROP TRIGGER staff_security ON staff;
DROP FUNCTION staff_security_change();
DROP TABLE security_settings, audit_events, refresh_tokens, accounts;
ALTER TABLE staff DROP COLUMN role;
