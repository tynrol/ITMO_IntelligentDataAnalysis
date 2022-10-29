CREATE TABLE Images (
    id TEXT PRIMARY KEY,
    width INT,
    height INT,
    description TEXT,
    url TEXT NOT NULL,
    path TEXT,
    created_at TIMESTAMP
);