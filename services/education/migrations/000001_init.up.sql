-- migrations/000001_init.up.sql
CREATE TABLE courses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    year INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    squad_id BIGINT NOT NULL,

    CONSTRAINT chk_ends_after_starts CHECK (ends_at > starts_at)
);

CREATE INDEX idx_courses_squad_id ON courses(squad_id);
CREATE INDEX idx_courses_squad_year ON courses(squad_id, year);
