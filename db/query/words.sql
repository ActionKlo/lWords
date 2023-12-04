-- name: ListWords :many
SELECT * FROM words;

-- name: CreateWord :one
INSERT INTO words (
    eng, rus, learn_at
) VALUES (
    $1, $2, CURRENT_TIMESTAMP
) ON CONFLICT (eng) DO UPDATE
SET rus = EXCLUDED.rus || '; ' || words.rus
RETURNING *;