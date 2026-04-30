-- migrations/000003_create_activities.up.sql
CREATE TABLE activities (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    max_score INTEGER NOT NULL DEFAULT 0 CHECK (max_score >= 0),
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    activity_type_id BIGINT NOT NULL REFERENCES activity_types(id),

    CONSTRAINT chk_ends_after_starts CHECK (ends_at > starts_at)
);

CREATE INDEX idx_activities_course_id ON activities(course_id);
CREATE INDEX idx_activities_type_id ON activities(activity_type_id);
