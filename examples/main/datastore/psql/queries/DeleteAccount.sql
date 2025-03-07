-- name: DeleteAccount :exec
-- Delete one Account using id
DELETE FROM accounts
WHERE id = @id ;