CREATE TABLE IF NOT EXISTS songs (
	id SERIAL,
	title VARCHAR(255) NOT NULL,
	artist VARCHAR(255) NOT NULL,
	song_url VARCHAR(255) NOT NULL,
	image_url VARCHAR(255),
    description VARCHAR(280),
	confirmed BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	confirmed_at TIMESTAMP WITH TIME ZONE,
    user_id INTEGER,
	UNIQUE (title),
	UNIQUE (song_url),
	PRIMARY KEY (id)
);
CREATE INDEX songs_confirmed_index ON songs (confirmed);
CREATE INDEX songs_created_at_index ON songs (created_at);
ALTER TABLE songs ADD CONSTRAINT songs_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);
INSERT INTO songs (id, title, artist, song_url, description, created_at) SELECT id, title, author, song_url, description, created_at FROM production_song;