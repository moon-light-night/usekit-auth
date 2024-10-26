INSERT INTO apps (id, name, secret)
VALUES (1, 'app', 'test-secret')
    ON CONFLICT DO NOTHING;