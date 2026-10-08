-- +goose Up
CREATE TABLE patients (
    id uuid PRIMARY KEY,
    first_name text NOT NULL CHECK (char_length(first_name) BETWEEN 1 AND 64),
    last_name text NOT NULL CHECK (char_length(last_name) BETWEEN 1 AND 64),
    middle_name text CHECK (middle_name IS NULL OR char_length(middle_name) BETWEEN 1 AND 64),
    birth_date date NOT NULL,
    administrative_sex text NOT NULL CHECK (administrative_sex IN ('M','F','UNKNOWN')),
    comment text NOT NULL DEFAULT '' CHECK (char_length(comment) <= 32),
    created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    updated_at timestamptz NOT NULL DEFAULT clock_timestamp()
);
CREATE TABLE records (
    id uuid PRIMARY KEY,
    patient_id uuid NOT NULL UNIQUE REFERENCES patients(id) ON DELETE RESTRICT,
    version bigint NOT NULL DEFAULT 1 CHECK (version >= 1),
    state jsonb NOT NULL CHECK (jsonb_typeof(state) = 'object'),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE record_events (
    id uuid PRIMARY KEY,
    record_id uuid NOT NULL REFERENCES records(id) ON DELETE RESTRICT,
    sequence bigint NOT NULL CHECK (sequence >= 1),
    version bigint NOT NULL CHECK (version >= 1),
    occurred_at timestamptz NOT NULL,
    actor jsonb NOT NULL CHECK (jsonb_typeof(actor) = 'object' AND actor ? 'id'),
    event_type text NOT NULL CHECK (event_type IN ('record_created','diagnosis_added','diagnosis_removed','prescription_added','prescription_cancelled','prescription_removed')),
    entity_id uuid NOT NULL,
    before_data jsonb NOT NULL,
    after_data jsonb NOT NULL,
    UNIQUE(record_id,sequence)
);
CREATE INDEX record_events_time_idx ON record_events(record_id,occurred_at,sequence);
CREATE INDEX record_events_version_idx ON record_events(record_id,version,sequence);

-- +goose Down
DROP TABLE record_events;
DROP TABLE records;
DROP TABLE patients;
