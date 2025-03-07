-- name: ListAccount :many
-- Lists 1000 Account having id > @id
SELECT * FROM accounts
WHERE id > @id
ORDER BY id
LIMIT 1000 ;