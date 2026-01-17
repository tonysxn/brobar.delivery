-- Add delivery door price setting
INSERT INTO settings (key, type, value) VALUES ('delivery_door_price', 'number', '50') ON CONFLICT DO NOTHING;
