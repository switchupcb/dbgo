-- name: InsertUser :one
-- Insert one row of User
INSERT INTO users
( 
    name
 ,  password
 ,  email
 ,  created_at
 ,  updated_at
) VALUES (
    @name
 ,  @password
 ,  @email
 ,  @created_at
 ,  @updated_at
)
RETURNING * ;