-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(100),
    email VARCHAR(100),
    role VARCHAR(50) CHECK (role IN ('student', 'teacher', 'admin')) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    last_logout TIMESTAMP
);

-- Добавление администратора
INSERT INTO users (telegram_id, name, email, role, active)
VALUES (6511775557, 'Admin User', 'admin@education-platform.com', 'admin', TRUE);

-- Таблица пропусков
CREATE TABLE IF NOT EXISTS passcodes (
    code VARCHAR(10) PRIMARY KEY,
    role VARCHAR(50) CHECK (role IN ('student', 'teacher', 'admin')) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP
);

-- Таблица предметов
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Добавление предметов
INSERT INTO subjects (name) VALUES
('Математика'),
('Информатика'),
('Физика'),
('Химия'),          -- Новый предмет
('Биология'),       -- Новый предмет
('История'),        -- Новый предмет
('Литература');     -- Новый предмет

-- Таблица групп
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Добавление групп
INSERT INTO groups (name) VALUES
('IT-21'),
('IT-22'),
('Math-101'),
('Bio-101'),        -- Новая группа
('Chem-202'),       -- Новая группа
('Hist-301');       -- Новая группа

-- Таблица преподавателей
CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Таблица расписания
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    subject_id INT NOT NULL,
    date DATE NOT NULL,
    time TIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (subject_id) REFERENCES subjects(id)
);