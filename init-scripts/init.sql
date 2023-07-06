CREATE USER edwin WITH PASSWORD 'acid';
CREATE DATABASE edwin;
GRANT ALL PRIVILEGES ON DATABASE edwin TO edwin;
\c edwin;
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    wallet INTEGER,
    permission BOOLEAN NOT NULL DEFAULT TRUE
);
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    admin_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE
);
INSERT INTO admins (created_at, updated_at, admin_name, email, phone, password, role, active)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'Edwin', 'admin01@gmail.com', '9048402133', '$2a$10$KV7
aw2CXyRZ1MeO0haWN2e1Q3FQTQaDQi/3AGT4XAN1LnhcrChoGq', 'master', TRUE);


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO edwin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO edwin;

