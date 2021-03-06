CREATE TABLE IF NOT EXISTS production_song (
	id SERIAL,
	title VARCHAR(255) NOT NULL,
	author VARCHAR(255) NOT NULL,
	song_url VARCHAR(255) NOT NULL,
	image_url VARCHAR(255),
	description VARCHAR(280),
	recommended BOOLEAN NOT NULL DEFAULT false,			
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	recommended_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	UNIQUE (title),
	UNIQUE (song_url),
	PRIMARY KEY (id)
);