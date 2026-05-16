-- migrations/000001_create_models.down.sql

-- Удаляем внешние ключи
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS fk_invitations_created_by;
ALTER TABLE invitations DROP CONSTRAINT IF EXISTS fk_invitations_squad_id;

-- Удаляем таблицы
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS education_institutions;
DROP TABLE IF EXISTS users;