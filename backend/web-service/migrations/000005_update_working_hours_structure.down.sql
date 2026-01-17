-- Revert to old working_hours structure
UPDATE settings SET value = '{
    "monday": {"start": "11:00", "end": "22:00", "closed": false},
    "tuesday": {"start": "11:00", "end": "22:00", "closed": false},
    "wednesday": {"start": "11:00", "end": "22:00", "closed": false},
    "thursday": {"start": "11:00", "end": "22:00", "closed": false},
    "friday": {"start": "11:00", "end": "23:00", "closed": false},
    "saturday": {"start": "11:00", "end": "23:00", "closed": false},
    "sunday": {"start": "11:00", "end": "22:00", "closed": false}
}' WHERE key = 'working_hours';
