CREATE TABLE roles
(
    id          CHAR(26) PRIMARY KEY,
    name        CITEXT      NOT NULL UNIQUE,
    created_at  TIMESTAMP   NOT NULL,
    modified_at TIMESTAMP,

    CONSTRAINT chk_roles_id_len CHECK (char_length(id) = 26),
    CONSTRAINT chk_roles_name_trimmed_nonempty CHECK (name = btrim(name) AND name <> '')
);

CREATE TABLE role_permissions
(
    role_id       CHAR(26) NOT NULL,
    permission_id CHAR(26) NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (id),
    FOREIGN KEY (permission_id) REFERENCES permissions (id)
);
