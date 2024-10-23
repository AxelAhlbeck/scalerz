-- name: GetAnswer :one
SELECT * FROM question
WHERE question = $1 LIMIT 1;

-- name: GetQuestion :one
SELECT * FROM question
WHERE id = $1 LIMIT 1;

-- name: InsertAnswer :one
INSERT INTO question (
    question, answer
) VALUES (
    $1, $2
)
RETURNING *;