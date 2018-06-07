CREATE TABLE IF NOT EXISTS users (
	id SERIAL,
	email VARCHAR(255) NOT NULL,
	username VARCHAR(25) NOT NULL,
	password_hash VARCHAR(128) NOT NULL,
	role_id INTEGER NOT NULL,
	UNIQUE (username),
	UNIQUE (email),
	PRIMARY KEY (id),
	CONSTRAINT users_role_id_fk FOREIGN KEY (role_id) REFERENCES roles (id)
);