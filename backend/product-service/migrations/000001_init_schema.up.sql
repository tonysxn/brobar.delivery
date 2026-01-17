-- Включение расширения для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица категорий
CREATE TABLE categories (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            name VARCHAR(255) NOT NULL,
                            slug VARCHAR(255) UNIQUE NOT NULL,
                            icon VARCHAR(255),
                            sort INTEGER DEFAULT 0
);

-- Таблица продуктов
CREATE TABLE products (
                          id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                          external_id VARCHAR(100) UNIQUE NOT NULL,
                          name VARCHAR(255) NOT NULL,
                          slug VARCHAR(255) UNIQUE NOT NULL,
                          description TEXT,
                          price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
                          weight DECIMAL(10,3) DEFAULT 0,
                          category_id uuid REFERENCES categories(id) ON DELETE CASCADE,
                          sort INTEGER DEFAULT 0,
                          hidden BOOLEAN DEFAULT FALSE,
                          alcohol BOOLEAN DEFAULT FALSE,
                          sold BOOLEAN DEFAULT FALSE,
                          image VARCHAR(255) NOT NULL
);

-- Индексы для ускорения поиска
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_slug ON products(slug);