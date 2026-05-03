-- Create initial tables for Squad & Membership Service
CREATE TABLE IF NOT EXISTS squads (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    short_name TEXT,
    region TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS squad_memberships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    squad_id BIGINT NOT NULL REFERENCES squads(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('candidate', 'fighter')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS squad_leaderships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    squad_id BIGINT NOT NULL REFERENCES squads(id) ON DELETE CASCADE,
    position TEXT NOT NULL CHECK (position IN ('commander', 'commissar', 'methodist', 'press_center_head', 'commandant')),
    appointed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dismissed_at TIMESTAMP
);
