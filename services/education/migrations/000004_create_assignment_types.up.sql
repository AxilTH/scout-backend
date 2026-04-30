CREATE TABLE assignment_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Предзаполненные типы заданий
INSERT INTO assignment_types (name, description) VALUES
('Домашняя работа', 'Регулярное домашнее задание'),
('Проект', 'Проектная работа'),
('Эссе', 'Письменная работа'),
('Тест', 'Автоматизированный тест');
