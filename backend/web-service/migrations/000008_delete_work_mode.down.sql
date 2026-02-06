INSERT INTO settings (key, type, value) VALUES ('work_mode', 'string', 'delivery') ON CONFLICT DO NOTHING;
