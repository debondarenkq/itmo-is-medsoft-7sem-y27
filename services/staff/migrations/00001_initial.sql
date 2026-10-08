-- +goose Up
CREATE TABLE staff (
 id uuid PRIMARY KEY,
 first_name text NOT NULL CHECK (char_length(first_name) BETWEEN 3 AND 24),
 last_name text NOT NULL CHECK (char_length(last_name) BETWEEN 3 AND 24),
 position text NOT NULL CHECK (char_length(position) BETWEEN 3 AND 24),
 created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 updated_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 deleted_at timestamptz
);

-- +goose Down
DROP TABLE staff;
