INSERT INTO apps (id, name, secret)
VALUES (2, 'test-2', 'test_secret_2')
ON CONFLICT DO NOTHING;