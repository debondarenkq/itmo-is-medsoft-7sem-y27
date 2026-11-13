-- +goose Up
CREATE TABLE connection(id boolean PRIMARY KEY DEFAULT true CHECK(id),session jsonb NOT NULL);
CREATE TABLE results(id uuid PRIMARY KEY,request_id uuid NOT NULL UNIQUE,created_at timestamptz NOT NULL DEFAULT clock_timestamp(),snapshot jsonb NOT NULL,delivered boolean NOT NULL DEFAULT false);
CREATE INDEX results_date ON results(created_at,id);
-- +goose Down
DROP TABLE results,connection;
