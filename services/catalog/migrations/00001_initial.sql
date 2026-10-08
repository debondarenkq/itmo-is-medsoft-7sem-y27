-- +goose Up
CREATE TABLE diagnoses (
 id uuid PRIMARY KEY,
 code text NOT NULL CHECK (char_length(code) BETWEEN 6 AND 6) UNIQUE,
 name text NOT NULL CHECK (char_length(name) BETWEEN 3 AND 24) UNIQUE,
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 deleted_at timestamptz
);

-- +goose Down
DROP TABLE diagnoses;
