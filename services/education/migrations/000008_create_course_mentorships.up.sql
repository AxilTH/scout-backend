CREATE TABLE course_mentorships (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    mentor_user_id BIGINT NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(course_id, mentor_user_id)
);

CREATE INDEX idx_course_mentorships_course ON course_mentorships(course_id);
CREATE INDEX idx_course_mentorships_mentor ON course_mentorships(mentor_user_id);