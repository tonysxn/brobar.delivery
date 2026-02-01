ALTER TABLE orders ADD COLUMN invoice_id VARCHAR(64);
CREATE INDEX idx_orders_invoice_id ON orders(invoice_id);
