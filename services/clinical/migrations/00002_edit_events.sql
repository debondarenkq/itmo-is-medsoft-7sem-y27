-- +goose Up
ALTER TABLE record_events DROP CONSTRAINT record_events_event_type_check;
ALTER TABLE record_events ADD CONSTRAINT record_events_event_type_check
    CHECK (event_type IN ('record_created', 'diagnosis_added', 'diagnosis_updated', 'diagnosis_removed', 'prescription_added', 'prescription_updated', 'prescription_cancelled', 'prescription_removed'));

-- +goose Down
ALTER TABLE record_events DROP CONSTRAINT record_events_event_type_check;
ALTER TABLE record_events ADD CONSTRAINT record_events_event_type_check
    CHECK (event_type IN ('record_created', 'diagnosis_added', 'diagnosis_removed', 'prescription_added', 'prescription_cancelled', 'prescription_removed'));
