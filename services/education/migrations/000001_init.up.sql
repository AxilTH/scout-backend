-- migrations/000001_init.up.sql
CREATE TABLE courses (
    id UUID PRIMARY KEY,
    squad_id UUID NOT NULL,
    year INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    starts_at DATE NOT NULL,
    ends_at DATE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);