-- name: InsertAccount :one
-- Insert one row of Account
INSERT INTO accounts
( 
    first_name
 ,  last_name
 ,  email
 ,  created_at
 ,  updated_at
) VALUES (
    @first_name
 ,  @last_name
 ,  @email
 ,  @created_at
 ,  @updated_at
)
RETURNING * ;