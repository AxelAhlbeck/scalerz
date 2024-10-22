-- name: GetAnswer :one
SELECT * FROM question
WHERE question = $1 LIMIT 1;

-- name: GetQustion :one
SELECT * FROM question
WHERE id = $1 LIMIT 1;