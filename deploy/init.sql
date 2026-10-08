-- Only local demonstration credentials; each service owns a separate database.
CREATE ROLE staff_user LOGIN PASSWORD 'staff_local';
CREATE ROLE catalog_user LOGIN PASSWORD 'catalog_local';
CREATE ROLE clinical_user LOGIN PASSWORD 'clinical_local';
CREATE DATABASE staff OWNER staff_user;
CREATE DATABASE catalog OWNER catalog_user;
CREATE DATABASE clinical OWNER clinical_user;
REVOKE CONNECT ON DATABASE staff FROM PUBLIC;
REVOKE CONNECT ON DATABASE catalog FROM PUBLIC;
REVOKE CONNECT ON DATABASE clinical FROM PUBLIC;
GRANT CONNECT ON DATABASE staff TO staff_user;
GRANT CONNECT ON DATABASE catalog TO catalog_user;
GRANT CONNECT ON DATABASE clinical TO clinical_user;
