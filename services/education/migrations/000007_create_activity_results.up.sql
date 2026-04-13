CREATE TABLE activity_results (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    activity_id BIGINT NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'not_submitted',
    score INTEGER CHECK (score >= 0),
    feedback TEXT,
    submitted_at TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewer_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, activity_id)
);

CREATE INDEX idx_activity_results_user ON activity_results(user_id);
CREATE INDEX idx_activity_results_activity ON activity_results(activity_id);
CREATE INDEX idx_activity_results_status ON activity_results(status);
