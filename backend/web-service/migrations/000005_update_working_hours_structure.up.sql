-- Update working_hours to new structure with separate delivery and pickup hours
UPDATE settings SET value = '{
  "delivery": {
    "monday": {"start": "11:00", "end": "22:00", "closed": false},
    "tuesday": {"start": "11:00", "end": "22:00", "closed": false},
    "wednesday": {"start": "11:00", "end": "22:00", "closed": false},
    "thursday": {"start": "11:00", "end": "22:00", "closed": false},
    "friday": {"start": "11:00", "end": "23:00", "closed": false},
    "saturday": {"start": "11:00", "end": "23:00", "closed": false},
    "sunday": {"start": "11:00", "end": "22:00", "closed": false}
  },
  "pickup": {
    "monday": {"start": "11:00", "end": "22:00", "closed": false},
    "tuesday": {"start": "11:00", "end": "22:00", "closed": false},
    "wednesday": {"start": "11:00", "end": "22:00", "closed": false},
    "thursday": {"start": "11:00", "end": "22:00", "closed": false},
    "friday": {"start": "11:00", "end": "23:00", "closed": false},
    "saturday": {"start": "11:00", "end": "23:00", "closed": false},
    "sunday": {"start": "11:00", "end": "22:00", "closed": false}
  }
}' WHERE key = 'working_hours';
