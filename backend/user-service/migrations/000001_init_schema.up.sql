-- Включение расширения для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создать ENUM тип для ролей
CREATE TYPE user_role AS ENUM ('admin', 'user', 'moderator');

-- Создать таблицу пользователей
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       role_id user_role NOT NULL DEFAULT 'user',
                       email TEXT NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       name TEXT NOT NULL,
                       address TEXT,
                       address_coords TEXT,
                       phone TEXT,
                       promo_card TEXT
);
