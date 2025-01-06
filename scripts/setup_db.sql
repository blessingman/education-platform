-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(100),
    email VARCHAR(100),
    role VARCHAR(50) CHECK (role IN ('student', 'teacher', 'admin')) NOT NULL
);

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

-- Таблица групп
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Таблица преподавателей
CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

-- Таблица расписания
CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    teacher_id INTEGER REFERENCES teachers(id),
    subject VARCHAR(100) NOT NULL,
    time TIMESTAMP NOT NULL,
    location VARCHAR(100) NOT NULL
);

-- Таблица учебных материалов
CREATE TABLE IF NOT EXISTS materials (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    title VARCHAR(200) NOT NULL,
    file_url VARCHAR(200) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица объявлений
CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица обратной связи
CREATE TABLE IF NOT EXISTS feedback (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
