-- Возвращаем старую колонку
ALTER TABLE assignments 
ADD COLUMN assignment_type TEXT NOT NULL DEFAULT 'homework';

-- Удаляем новую колонку
ALTER TABLE assignments 
DROP COLUMN assignment_type_id;
