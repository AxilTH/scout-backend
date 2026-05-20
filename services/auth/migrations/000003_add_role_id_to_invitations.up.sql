-- Добавляем колонку role_id в таблицу приглашений
ALTER TABLE invitations
ADD COLUMN role_id BIGINT NOT NULL DEFAULT 1;

-- Обновляем внешние ключи если нужно (в данном случае нет внешнего ключа для role_id, 
-- но если бы был, добавили бы его здесь)
