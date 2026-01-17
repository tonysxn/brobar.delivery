INSERT INTO settings (key, type, value) VALUES ('working_hours', 'json', '{
    "monday": {"start": "11:00", "end": "22:00", "closed": false},
    "tuesday": {"start": "11:00", "end": "22:00", "closed": false},
    "wednesday": {"start": "11:00", "end": "22:00", "closed": false},
    "thursday": {"start": "11:00", "end": "22:00", "closed": false},
    "friday": {"start": "11:00", "end": "23:00", "closed": false},
    "saturday": {"start": "11:00", "end": "23:00", "closed": false},
    "sunday": {"start": "11:00", "end": "22:00", "closed": false}
}') ON CONFLICT DO NOTHING;

INSERT INTO settings (key, type, value) VALUES ('sales_paused', 'boolean', 'false') ON CONFLICT DO NOTHING;
