CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE If NOT EXISTS Expressions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expressionID TEXT NOT NULL UNIQUE,
    expression TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'processing', 'completed', 'failed')),
    result REAL,
    FOREIGN KEY (user_id) REFERENCES users (id)
)