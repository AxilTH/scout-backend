-- Добавляем новую колонку
ALTER TABLE assignments 
ADD COLUMN assignment_type_id BIGINT REFERENCES assignment_types(id);

-- Заполняем значениями по умолчанию (для существующих данных)
UPDATE assignments 
SET assignment_type_id = 1
WHERE assignment_type_id IS NULL;

-- Делаем колонку NOT NULL
ALTER TABLE assignments 
ALTER COLUMN assignment_type_id SET NOT NULL;

-- Удаляем старую колонку
ALTER TABLE assignments 
DROP COLUMN assignment_type;
