CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE permissions
(
    id          CHAR(26)  PRIMARY KEY,
    name        CITEXT    NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP NOT NULL,
    modified_at TIMESTAMP,
    CONSTRAINT chk_permissions_id_len CHECK (char_length(id) = 26),
    CONSTRAINT chk_permissions_name_trimmed_nonempty CHECK (name = btrim(name) AND name <> '')
);
