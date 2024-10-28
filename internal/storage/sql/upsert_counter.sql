INSERT INTO metrics (name, type, counter)
VALUES ($1, $2, $3)
ON CONFLICT (name, type)
DO UPDATE SET counter = metrics.counter + EXCLUDED.counter
