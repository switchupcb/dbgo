-- name: UpdateUser :one
-- Update one row of User using id
UPDATE users
SET 
    name = @name
 ,  password = @password
 ,  email = @email
 ,  created_at = @created_at
 ,  updated_at = @updated_at
WHERE id = @id
RETURNING * ;