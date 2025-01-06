-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    name TEXT,
    email TEXT UNIQUE,
    role TEXT CHECK (role IN ('student', 'teacher', 'admin')) NOT NULL
);

-- Таблица пропусков
CREATE TABLE passcodes (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    role TEXT CHECK (role IN ('student', 'teacher', 'admin')) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP
);

-- Таблица групп (пример)
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Таблица преподавателей (пример)
CREATE TABLE teachers (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Таблица расписаний (пример)
CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    teacher_id INTEGER REFERENCES teachers(id),
    subject TEXT NOT NULL,
    time TIME NOT NULL,
    location TEXT NOT NULL,
    day DATE NOT NULL
);

-- Таблица учебных материалов (пример)
CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    title TEXT NOT NULL,
    file_url TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица обратной связи (пример)
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица объявлений (пример)
CREATE TABLE announcements (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
