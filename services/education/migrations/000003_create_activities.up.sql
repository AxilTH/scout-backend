CREATE TABLE activities (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    activity_type_id UUID NOT NULL REFERENCES activity_types(id),
    title TEXT NOT NULL,
    description TEXT,
    starts_at TIMESTAMP NOT NULL,
    ends_at TIMESTAMP NOT NULL,
    max_score INTEGER CHECK (max_score >= 0),
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_activities_course_id ON activities(course_id);
CREATE INDEX idx_activities_type_id ON activities(activity_type_id);