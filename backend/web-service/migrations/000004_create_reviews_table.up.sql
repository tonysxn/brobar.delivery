CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_rating INT NOT NULL CHECK (food_rating >= 1 AND food_rating <= 5),
    service_rating INT NOT NULL CHECK (service_rating >= 1 AND service_rating <= 5),
    comment TEXT DEFAULT '',
    phone VARCHAR(20),
    email VARCHAR(255),
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC);
