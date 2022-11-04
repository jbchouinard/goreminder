CREATE USER mxremind WITH PASSWORD 'mxremind';

CREATE DATABASE mxremind OWNER mxremind;

\c mxremind;

CREATE TABLE reminders (
    id UUID PRIMARY KEY,
    generated_from_id TEXT,
    recipient TEXT,
    content TEXT,
    due_time TIMESTAMP,
    is_sent BOOLEAN
);

ALTER TABLE public.reminders OWNER TO mxremind;
