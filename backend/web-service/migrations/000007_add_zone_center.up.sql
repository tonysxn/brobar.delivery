-- Add zone_center setting
INSERT INTO settings (key, value, setting_type) VALUES ('zone_center', '{"lat": 50.0014656, "lng": 36.245192}', 'json')
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;
