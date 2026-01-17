CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    value TEXT NOT NULL
);

INSERT INTO settings (key, type, value) VALUES ('work_mode', 'string', 'delivery') ON CONFLICT DO NOTHING;
INSERT INTO settings (key, type, value) VALUES ('delivery_radius', 'number', '10000') ON CONFLICT DO NOTHING;
INSERT INTO settings (key, type, value) VALUES ('delivery_zones', 'json', '[
    {"radius": 2, "innerRadius": 0, "color": "#22c55e", "price": 150, "freeOrderPrice": 600, "name": "Зелена зона"},
    {"radius": 5, "innerRadius": 2, "color": "#eab308", "price": 200, "freeOrderPrice": 1100, "name": "Жовта зона"},
    {"radius": 7, "innerRadius": 5, "color": "#f97316", "price": 300, "freeOrderPrice": 1800, "name": "Помаранчева зона"},
    {"radius": 10, "innerRadius": 7, "color": "#ef4444", "price": 350, "freeOrderPrice": 2400, "name": "Червона зона"}
]') ON CONFLICT DO NOTHING;
