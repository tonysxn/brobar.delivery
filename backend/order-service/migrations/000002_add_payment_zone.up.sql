-- Add payment_method and zone columns to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS payment_method VARCHAR(32) DEFAULT 'cash';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS zone VARCHAR(64);
