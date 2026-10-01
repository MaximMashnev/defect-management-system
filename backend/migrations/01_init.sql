-- 01_init.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'manager',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS buildings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    date_registration TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    number_applications INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    score INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    email VARCHAR(255) NOT NULL,
    date_submission TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    date_in_progress TIMESTAMPTZ,
    date_result TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL DEFAULT 'новый',
    priority INT NOT NULL DEFAULT 1,
    building_id INT NOT NULL REFERENCES buildings(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS uploads (
    id SERIAL PRIMARY KEY,
    application_id INT REFERENCES applications(id) ON DELETE SET NULL,
    original_name VARCHAR(255) NOT NULL,
    stored_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    url TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS application_uploads (
    application_id INT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    upload_id INT NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    PRIMARY KEY (application_id, upload_id)
);

-- Seed initial data
-- Password for admin@a.a is admin123 (bcrypt hash: $2a$10$3mEvG99sW7XG7Kx0c/lZ7uXQ5t0QW1ZkYnK1J3p9z0j/1G5B2sF2C)
INSERT INTO users (login, password_hash, role)
VALUES ('admin@a.a', '$2a$10$63kR9wS8v5b9vR4wF.e.Teu9Kk3zK7kK7Z7HqPz3o8M8J3l6d8qfa', 'manager')
ON CONFLICT (login) DO NOTHING;

-- Initial categories matching test-task-react
INSERT INTO categories (id, name, score) VALUES
(1, 'Водоснабжение и отопление', 3),
(2, 'Электрика', 3),
(3, 'Окна и двери', 2),
(4, 'Вентиляция', 2),
(5, 'Отделочные работы', 1),
(6, 'Фасад и кровля', 2)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, score = EXCLUDED.score;

-- Initial buildings
INSERT INTO buildings (id, name, address, date_registration, number_applications) VALUES
(1, 'ЖК "Северное Сияние", Корпус 1', 'г. Москва, ул. Полярная, д. 14', NOW() - INTERVAL '30 days', 3),
(2, 'ЖК "Северное Сияние", Корпус 2', 'г. Москва, ул. Полярная, д. 16', NOW() - INTERVAL '25 days', 2),
(3, 'ЖК "Лазурный Берег", Башня А', 'г. Санкт-Петербург, Приморский пр., д. 45', NOW() - INTERVAL '60 days', 4),
(4, 'БЦ "Горизонт"', 'г. Москва, Ленинградский пр-т, д. 80', NOW() - INTERVAL '15 days', 1)
ON CONFLICT (id) DO NOTHING;

SELECT setval('categories_id_seq', (SELECT COALESCE(MAX(id), 1) FROM categories));
SELECT setval('buildings_id_seq', (SELECT COALESCE(MAX(id), 1) FROM buildings));
