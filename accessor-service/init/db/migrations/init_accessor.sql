CREATE TABLE users (
   session_id       TEXT PRIMARY KEY,
   lied             INTEGER
);

CREATE TABLE images (
    id              TEXT,
    session_id      TEXT,
    type            TEXT,
    width           INT,
    height          INT,
    description     TEXT,
    url             TEXT NOT NULL,
    path            TEXT,
    created_at      TIMESTAMP,
    PRIMARY KEY(id, session_id),
    FOREIGN KEY(session_id) REFERENCES users
);
