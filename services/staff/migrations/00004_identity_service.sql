-- +goose Up
-- Refuse to discard an existing installation before the explicit transfer.
-- +goose StatementBegin
DO $$ BEGIN
 IF (EXISTS(SELECT 1 FROM accounts) OR EXISTS(SELECT 1 FROM audit_events) OR EXISTS(SELECT 1 FROM security_settings))
    AND to_regclass('_identity_transfer') IS NULL THEN
  RAISE EXCEPTION 'Run make migrate-identity before starting staff-service; existing identity data must be preserved';
 END IF;
END $$;
-- +goose StatementEnd
DROP TRIGGER staff_security ON staff;
DROP FUNCTION staff_security_change();
DROP TABLE security_settings,audit_events,refresh_tokens,accounts;
CREATE SEQUENCE staff_identity_revision;
ALTER TABLE staff ADD COLUMN identity_revision bigint NOT NULL DEFAULT nextval('staff_identity_revision');
ALTER TABLE staff ADD COLUMN security_revision bigint NOT NULL DEFAULT 0;
CREATE TABLE identity_outbox (
 revision bigint PRIMARY KEY,
 staff_id uuid NOT NULL,
 payload jsonb NOT NULL
);
INSERT INTO identity_outbox SELECT identity_revision,id,to_jsonb(staff) FROM staff;
-- +goose StatementBegin
CREATE FUNCTION staff_identity_version() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 NEW.identity_revision=nextval('staff_identity_revision');
 IF TG_OP='UPDATE' AND (NEW.role IS DISTINCT FROM OLD.role OR (NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL)) THEN
  NEW.security_revision=NEW.identity_revision;
 END IF;
 RETURN NEW;
END $$;
-- +goose StatementEnd
CREATE TRIGGER staff_identity_version BEFORE INSERT OR UPDATE ON staff FOR EACH ROW EXECUTE FUNCTION staff_identity_version();
-- +goose StatementBegin
CREATE FUNCTION staff_identity_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 INSERT INTO identity_outbox VALUES(NEW.identity_revision,NEW.id,to_jsonb(NEW));
 RETURN NEW;
END $$;
-- +goose StatementEnd
CREATE TRIGGER staff_identity_event AFTER INSERT OR UPDATE ON staff FOR EACH ROW EXECUTE FUNCTION staff_identity_event();
-- +goose Down
-- This boundary change requires restoring the pre-upgrade backup, not dropping identity data.
-- +goose StatementBegin
DO $$ BEGIN RAISE EXCEPTION 'Restore a pre-identity backup to roll back this service split'; END $$;
-- +goose StatementEnd
