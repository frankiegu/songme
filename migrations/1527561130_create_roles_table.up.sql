CREATE TABLE IF NOT EXISTS roles (
	id SERIAL,
	name VARCHAR(64) NOT NULL,
    default_role BOOLEAN NOT NULL DEFAULT false,
    permissions INTEGER NOT NULL,
	UNIQUE (name),
	PRIMARY KEY (id)
);
CREATE INDEX roles_default_column_index ON roles (default_role);
INSERT INTO roles (name, default_role, permissions) VALUES ('User', true, 1);
INSERT INTO roles (name, default_role, permissions) VALUES ('Moderator', false, 3);
INSERT INTO roles (name, default_role, permissions) VALUES ('Administrator', false, 7);