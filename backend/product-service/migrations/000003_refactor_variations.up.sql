-- 000003_refactor_variations.up.sql

-- Remove global unique constraints
ALTER TABLE product_variations DROP CONSTRAINT IF EXISTS product_variations_external_id_key;
ALTER TABLE product_variation_groups DROP CONSTRAINT IF EXISTS product_variation_groups_external_id_key;

-- Add composite unique constraints
ALTER TABLE product_variations ADD CONSTRAINT uq_product_variations_group_external_id UNIQUE (group_id, external_id);
ALTER TABLE product_variation_groups ADD CONSTRAINT uq_product_variation_groups_product_external_id UNIQUE (product_id, external_id);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_product_variations_group_id ON product_variations(group_id);
CREATE INDEX IF NOT EXISTS idx_product_variation_groups_product_id ON product_variation_groups(product_id);
