-- +goose Up
ALTER TABLE accounts ADD COLUMN is_service boolean NOT NULL DEFAULT false;
ALTER TABLE accounts ADD CONSTRAINT service_account_kind CHECK(NOT is_service OR (NOT is_admin AND staff_id IS NULL));
CREATE UNIQUE INDEX service_key_hash ON accounts(password_hash) WHERE is_service;
-- +goose Down
DROP INDEX service_key_hash;
ALTER TABLE accounts DROP CONSTRAINT service_account_kind;
ALTER TABLE accounts DROP COLUMN is_service;
