-- Drop tables in reverse order
DROP TRIGGER IF EXISTS update_regions_updated_at ON regions;
DROP TRIGGER IF EXISTS update_positions_updated_at ON positions;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_squad_leaderships_updated_at ON squad_leaderships;
DROP TRIGGER IF EXISTS update_squad_memberships_updated_at ON squad_memberships;
DROP TRIGGER IF EXISTS update_squads_updated_at ON squads;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS squad_leaderships;
DROP TABLE IF EXISTS squad_memberships;
DROP TABLE IF EXISTS squads;
DROP TABLE IF EXISTS regions;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS roles;
