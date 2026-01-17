-- 000003_refactor_variations.down.sql

-- Remove composite unique constraints
ALTER TABLE product_variations DROP CONSTRAINT IF EXISTS uq_product_variations_group_external_id;
ALTER TABLE product_variation_groups DROP CONSTRAINT IF EXISTS uq_product_variation_groups_product_external_id;

-- Restore global unique constraints
-- Note: This might fail if there are duplicates, which is expected in dev
ALTER TABLE product_variations ADD CONSTRAINT product_variations_external_id_key UNIQUE (external_id);
ALTER TABLE product_variation_groups ADD CONSTRAINT product_variation_groups_external_id_key UNIQUE (external_id);
