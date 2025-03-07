-- name: ListUser :many
-- Lists 1000 User having id > @id
SELECT * FROM users
WHERE id > @id
ORDER BY id
LIMIT 1000 ;