DROP INDEX IF EXISTS idx_orders_invoice_id;
ALTER TABLE orders DROP COLUMN IF EXISTS invoice_id;
