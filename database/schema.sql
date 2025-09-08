-- Создание базы данных и таблиц для доски задач

-- Таблица досок
CREATE TABLE boards (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица колонок (фиксированные три колонки)
CREATE TABLE columns (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    order_num INTEGER NOT NULL
);

-- Вставляем три предустановленные колонки
INSERT INTO columns (id, name, order_num) VALUES
    ('todo', 'Актуальные задачи', 1),
    ('in-progress', 'В работе', 2),
    ('done', 'Выполнено', 3);

-- Таблица карточек
CREATE TABLE cards (
    id VARCHAR(32) PRIMARY KEY,
    board_id VARCHAR(32) NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    column_id VARCHAR(20) NOT NULL REFERENCES columns(id),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    assignee VARCHAR(255),
    order_num INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX idx_cards_board_id ON cards(board_id);
CREATE INDEX idx_cards_column_id ON cards(column_id);
CREATE INDEX idx_cards_board_column ON cards(board_id, column_id);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для автоматического обновления updated_at в таблице cards
CREATE TRIGGER update_cards_updated_at
    BEFORE UPDATE ON cards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
