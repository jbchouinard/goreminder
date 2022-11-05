CREATE TABLE reminders (
    id UUID PRIMARY KEY,
    generated_from_id TEXT,
    recipient TEXT,
    content TEXT,
    due_time TIMESTAMP,
    is_sent BOOLEAN
);

CREATE TABLE users (
    email TEXT PRIMARY KEY,
    timezone TEXT
);

---- create above / drop below ----

DROP TABLE users;

DROP TABLE reminders;
