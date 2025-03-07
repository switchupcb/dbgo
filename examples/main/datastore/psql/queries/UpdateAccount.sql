-- name: UpdateAccount :one
-- Update one row of Account using id
UPDATE accounts
SET 
    first_name = @first_name
 ,  last_name = @last_name
 ,  email = @email
 ,  created_at = @created_at
 ,  updated_at = @updated_at
WHERE id = @id
RETURNING * ;