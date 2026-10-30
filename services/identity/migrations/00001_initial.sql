-- +goose Up
CREATE TABLE staff_profiles (
 id uuid PRIMARY KEY,
 first_name text NOT NULL,
 last_name text NOT NULL,
 role text NOT NULL CHECK(role IN ('doctor','nurse','laboratory-technician','other')),
 deleted_at timestamptz,
 security_revision bigint NOT NULL DEFAULT 0,
 revision bigint NOT NULL
);
CREATE TABLE accounts (
 id uuid PRIMARY KEY,
 username text NOT NULL UNIQUE CHECK(char_length(username) BETWEEN 3 AND 16),
 password_hash text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 is_active boolean NOT NULL DEFAULT true,
 staff_id uuid UNIQUE REFERENCES staff_profiles(id),
 is_admin boolean NOT NULL DEFAULT false,
 totp_secret text,
 last_totp_step bigint NOT NULL DEFAULT -1,
 CHECK((is_admin AND staff_id IS NULL AND username='admin') OR NOT is_admin)
);
CREATE UNIQUE INDEX one_admin ON accounts(is_admin) WHERE is_admin;
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
 message text NOT NULL,
 event_type text NOT NULL DEFAULT 'other'
);
CREATE INDEX audit_time ON audit_events(event_time,event_id);
CREATE INDEX audit_event_type_time ON audit_events(event_type,event_time,event_id);
CREATE TABLE security_settings(id boolean PRIMARY KEY DEFAULT true CHECK(id),setup_secret text NOT NULL);
-- +goose StatementBegin
CREATE FUNCTION profile_security_change() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE affected uuid; revoked integer;
BEGIN
 IF NEW.security_revision > OLD.security_revision OR NEW.role IS DISTINCT FROM OLD.role OR (NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL) THEN
  SELECT id INTO affected FROM accounts WHERE staff_id=NEW.id FOR UPDATE;
  IF affected IS NOT NULL THEN
   UPDATE accounts SET is_active=CASE WHEN NEW.deleted_at IS NOT NULL THEN false ELSE is_active END,
    updated_at=clock_timestamp() WHERE id=affected;
   DELETE FROM refresh_tokens WHERE account_id=affected;
   GET DIAGNOSTICS revoked = ROW_COUNT;
   INSERT INTO audit_events(event_id,actor,message,event_type)
    VALUES(gen_random_uuid(),'system:staff:'||NEW.id,
    'Изменение роли/удаление медработника; аккаунт '||affected||'; отозвано refresh токенов: '||revoked,'account_updated');
  END IF;
 END IF;
 RETURN NEW;
END $$;
-- +goose StatementEnd
CREATE TRIGGER profile_security AFTER UPDATE ON staff_profiles FOR EACH ROW EXECUTE FUNCTION profile_security_change();
CREATE TABLE legacy_import(id boolean PRIMARY KEY DEFAULT true CHECK(id), source_hash text NOT NULL);
-- +goose Down
DROP TABLE legacy_import, security_settings,audit_events,refresh_tokens,accounts,staff_profiles;
DROP FUNCTION profile_security_change();
