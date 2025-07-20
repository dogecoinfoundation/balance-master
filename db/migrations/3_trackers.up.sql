 create table if not exists trackers (
    id INTEGER PRIMARY KEY,
	address TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS trackers_key ON trackers (id);