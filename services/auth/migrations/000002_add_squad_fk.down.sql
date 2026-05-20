-- Удаляем внешний ключ на таблицу squads
ALTER TABLE invitations
   DROP CONSTRAINT IF EXISTS fk_invitations_squad_id;