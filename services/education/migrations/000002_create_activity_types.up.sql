-- migrations/000002_create_activity_types.up.sql
CREATE TABLE activity_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    squad_id BIGINT,

    CONSTRAINT chk_squad_id_positive CHECK (squad_id IS NULL OR squad_id > 0),
    CONSTRAINT uniq_activity_type_title UNIQUE NULLS NOT DISTINCT (title, squad_id)
);

CREATE INDEX idx_activity_types_squad_id ON activity_types(squad_id);
CREATE INDEX idx_activity_types_title ON activity_types(title);

-- Сидирование глобальных типов (доступны всем отрядам)
INSERT INTO activity_types (title, description, squad_id) VALUES
    ('Лекция', 'Теоретическое занятие с презентацией материала', NULL),
    ('Домашнее задание', 'Самостоятельная работа вне занятий', NULL),
    ('Методический выезд', 'Выездное методическое мероприятие', NULL),
    ('Зачёт', 'Промежуточная форма контроля знаний', NULL),
    ('Экзамен', 'Итоговая аттестация по курсу', NULL);
