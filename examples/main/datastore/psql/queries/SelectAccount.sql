
-- name: SelectAccount :one
-- Select one Account using id
SELECT
    id
 ,  first_name
 ,  last_name
 ,  email
 ,  created_at
 ,  updated_at
FROM accounts
WHERE id = @id ;