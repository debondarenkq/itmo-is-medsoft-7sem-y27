-- +goose Up
ALTER TABLE records ADD COLUMN number bigint GENERATED ALWAYS AS IDENTITY;
ALTER TABLE records ADD CONSTRAINT records_number_unique UNIQUE(number);

-- +goose Down
ALTER TABLE records DROP COLUMN number;
