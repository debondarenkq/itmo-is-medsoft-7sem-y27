-- +goose Up
CREATE TABLE accounts (
 id uuid PRIMARY KEY,
 username text NOT NULL UNIQUE CHECK(char_length(username) BETWEEN 3 AND 16),
 password_hash text NOT NULL,
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 is_active boolean NOT NULL DEFAULT true,
 staff_id uuid UNIQUE,
 session_version bigint NOT NULL DEFAULT 0,
 is_admin boolean NOT NULL DEFAULT false,
 totp_secret text,
 last_totp_step bigint NOT NULL DEFAULT -1,
 CHECK((is_admin AND staff_id IS NULL AND username='admin') OR NOT is_admin)
);
CREATE UNIQUE INDEX one_admin ON accounts(is_admin) WHERE is_admin;
CREATE TABLE refresh_tokens (
 token_hash text PRIMARY KEY,
 account_id uuid NOT NULL REFERENCES accounts(id),
 expires_at timestamptz NOT NULL,
 staff_revision bigint NOT NULL DEFAULT 0
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
-- +goose Down
DROP TABLE security_settings,audit_events,refresh_tokens,accounts;
