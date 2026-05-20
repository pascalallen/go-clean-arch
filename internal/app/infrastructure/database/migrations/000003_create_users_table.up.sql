CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id          CHAR(26)        NOT NULL,
    first_name  VARCHAR(100)    NOT NULL CHECK (btrim(first_name) <> ''),
    last_name   VARCHAR(100)    NOT NULL CHECK (btrim(last_name) <> ''),
    email_address CITEXT        NOT NULL CHECK (btrim(email_address::text) <> ''),
    password_hash VARCHAR(255),
    created_at  TIMESTAMP       NOT NULL,
    modified_at TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_id_check CHECK (id ~ '^[0-9A-Z]{26}$'),
    CONSTRAINT users_email_address_key UNIQUE (email_address)
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id CHAR(26) NOT NULL REFERENCES users(id),
    role_id CHAR(26) NOT NULL REFERENCES roles(id),
    CONSTRAINT user_roles_pkey PRIMARY KEY (user_id, role_id)
);
