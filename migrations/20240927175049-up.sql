
-- +migrate Up
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    s_group VARCHAR(125) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date VARCHAR(55) DEFAULT '',
    text TEXT DEFAULT '',
    link VARCHAR(255) DEFAULT '',
    CONSTRAINT check_song_not_empty CHECK (song <> '')
    );

-- +migrate Down
