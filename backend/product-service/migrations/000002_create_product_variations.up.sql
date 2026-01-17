-- up.sql

CREATE TABLE product_variation_groups (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                          product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                                          name VARCHAR(255) NOT NULL,
                                          external_id VARCHAR(100) UNIQUE NOT NULL,
                                          default_value INT NULL,
                                          show BOOLEAN DEFAULT TRUE,
                                          required BOOLEAN DEFAULT FALSE
);

CREATE TABLE product_variations (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    group_id UUID NOT NULL REFERENCES product_variation_groups(id) ON DELETE CASCADE,
                                    external_id VARCHAR(100) UNIQUE NOT NULL,
                                    default_value INT NULL,
                                    show BOOLEAN DEFAULT TRUE,
                                    name VARCHAR(255) NOT NULL
);

CREATE INDEX idx_product_variation_groups_product_id ON product_variation_groups(product_id);
CREATE INDEX idx_product_variations_group_id ON product_variations(group_id);
