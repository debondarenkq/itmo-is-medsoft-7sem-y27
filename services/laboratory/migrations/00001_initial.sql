-- +goose Up
CREATE TABLE lab_tests(id uuid PRIMARY KEY,code text NOT NULL UNIQUE CHECK(code ~ '^\d{1,8}-\d$'),name text NOT NULL CHECK(char_length(name) BETWEEN 3 AND 160),active boolean NOT NULL DEFAULT true);
INSERT INTO lab_tests(id,code,name) VALUES
('30000000-0000-4000-8000-000000000001','2345-7','Глюкоза в крови'),
('30000000-0000-4000-8000-000000000002','33706-3','Антитела к Yersinia pestis (чума)'),
('30000000-0000-4000-8000-000000000003','718-7','Гемоглобин'),
('30000000-0000-4000-8000-000000000004','56888-1','ВИЧ 1/2: антитела и антиген p24'),
('30000000-0000-4000-8000-000000000005','5195-3','Гепатит B: HBsAg'),
('30000000-0000-4000-8000-000000000006','13955-0','Гепатит C: антитела'),
('30000000-0000-4000-8000-000000000007','22587-0','Сифилис: антитела к Treponema pallidum');
CREATE TABLE lab_requests(id uuid PRIMARY KEY,patient_id uuid NOT NULL,created_at timestamptz NOT NULL DEFAULT clock_timestamp(),resource jsonb NOT NULL);
CREATE INDEX lab_requests_patient ON lab_requests(patient_id,created_at,id);
CREATE TABLE lab_observations(id uuid PRIMARY KEY,request_id uuid NOT NULL UNIQUE REFERENCES lab_requests(id),patient_id uuid NOT NULL,created_at timestamptz NOT NULL DEFAULT clock_timestamp(),resource jsonb NOT NULL);
CREATE INDEX lab_observations_patient ON lab_observations(patient_id,created_at,id);
-- +goose Down
DROP TABLE lab_observations,lab_requests,lab_tests;
