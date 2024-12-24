INSERT INTO metrics (name, type, value)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET value = EXCLUDED.value
