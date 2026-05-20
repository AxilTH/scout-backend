-- Убираем колонку role_id из таблицы приглашений
ALTER TABLE invitations
DROP COLUMN role_id;
