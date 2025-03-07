-- name: DeleteUser :exec
-- Delete one User using id
DELETE FROM users
WHERE id = @id ;