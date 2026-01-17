CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE order_status AS ENUM (
    'pending',
    'paid',
    'shipping',
    'completed',
    'cancelled'
);

CREATE TYPE delivery_type AS ENUM (
    'delivery',
    'pickup',
    'dine'
);

CREATE TABLE orders (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NULL,
                        status_id order_status NOT NULL DEFAULT 'pending',
                        total_price NUMERIC(10,2) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
                        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

                        address VARCHAR(256) NOT NULL,
                        entrance VARCHAR(100),
                        floor VARCHAR(10),
                        flat VARCHAR(10),
                        address_wishes VARCHAR(1024),
                        name VARCHAR(100) NOT NULL,
                        phone VARCHAR(32),
                        time TIMESTAMP WITH TIME ZONE NOT NULL,
                        email VARCHAR(64),
                        wishes TEXT,
                        promo VARCHAR(128),
                        coords VARCHAR(256),
                        cutlery INTEGER,
                        delivery_cost NUMERIC(10,2),
                        delivery_door BOOLEAN DEFAULT FALSE,
                        delivery_type_id delivery_type NOT NULL DEFAULT 'delivery'
);

CREATE TABLE order_items (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             product_id UUID NOT NULL,
                             external_product_id VARCHAR(100) NOT NULL,
                             name VARCHAR(100) NOT NULL,
                             price NUMERIC(10,2) NOT NULL,
                             quantity INTEGER NOT NULL,
                             total_price NUMERIC(10,2) NOT NULL,
                             weight NUMERIC(10,3),
                             total_weight NUMERIC(10,3),

                             product_variation_group_id UUID,
                             product_variation_group_name VARCHAR(255),
                             product_variation_id UUID,
                             product_variation_external_id VARCHAR(100),
                             product_variation_name VARCHAR(255)
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
