-- name: SelectUser :one
-- Select one User using id
SELECT
    id
 ,  name
 ,  password
 ,  email
 ,  created_at
 ,  updated_at
FROM users
WHERE id = @id ;