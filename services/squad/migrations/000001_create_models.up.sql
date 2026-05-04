-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create positions table
CREATE TABLE IF NOT EXISTS positions (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create regions table
CREATE TABLE IF NOT EXISTS regions (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create squads table
CREATE TABLE IF NOT EXISTS squads (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    region_id BIGINT NOT NULL REFERENCES regions(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for squads
CREATE INDEX IF NOT EXISTS idx_squads_region_id ON squads(region_id);
CREATE INDEX IF NOT EXISTS idx_squads_created_at ON squads(created_at);

-- Create squad_memberships table
CREATE TABLE IF NOT EXISTS squad_memberships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    squad_id BIGINT NOT NULL REFERENCES squads(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Create indexes for squad_memberships
CREATE INDEX IF NOT EXISTS idx_squad_memberships_user_id ON squad_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_squad_memberships_squad_id ON squad_memberships(squad_id);
CREATE INDEX IF NOT EXISTS idx_squad_memberships_role_id ON squad_memberships(role_id);
CREATE INDEX IF NOT EXISTS idx_squad_memberships_is_active ON squad_memberships(is_active);
CREATE INDEX IF NOT EXISTS idx_squad_memberships_joined_at ON squad_memberships(joined_at);

-- Create squad_leaderships table
CREATE TABLE IF NOT EXISTS squad_leaderships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    squad_id BIGINT NOT NULL REFERENCES squads(id) ON DELETE CASCADE,
    position_id BIGINT NOT NULL REFERENCES positions(id) ON DELETE RESTRICT,
    appointed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    dismissed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for squad_leaderships
CREATE INDEX IF NOT EXISTS idx_squad_leaderships_user_id ON squad_leaderships(user_id);
CREATE INDEX IF NOT EXISTS idx_squad_leaderships_squad_id ON squad_leaderships(squad_id);
CREATE INDEX IF NOT EXISTS idx_squad_leaderships_position_id ON squad_leaderships(position_id);
CREATE INDEX IF NOT EXISTS idx_squad_leaderships_appointed_at ON squad_leaderships(appointed_at);
CREATE INDEX IF NOT EXISTS idx_squad_leaderships_dismissed_at ON squad_leaderships(dismissed_at);

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_squads_updated_at BEFORE UPDATE ON squads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_squad_memberships_updated_at BEFORE UPDATE ON squad_memberships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_squad_leaderships_updated_at BEFORE UPDATE ON squad_leaderships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_positions_updated_at BEFORE UPDATE ON positions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_regions_updated_at BEFORE UPDATE ON regions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
