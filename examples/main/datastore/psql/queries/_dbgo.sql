-- name: CountAccount :one
-- Count # of Account
SELECT count(*) as account_count from accounts ;

-- name: CountUser :one
-- Count # of User
SELECT count(*) as user_count from users ;

-- name: DeleteAccount :exec
-- Delete one Account using id
DELETE FROM accounts
WHERE id = @id ;

-- name: DeleteUser :exec
-- Delete one User using id
DELETE FROM users
WHERE id = @id ;

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

-- name: ListAccount :many
-- Lists 1000 Account having id > @id
SELECT * FROM accounts
WHERE id > @id
ORDER BY id
LIMIT 1000 ;

-- name: ListUser :many
-- Lists 1000 User having id > @id
SELECT * FROM users
WHERE id > @id
ORDER BY id
LIMIT 1000 ;


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

-- name: example :exec
SELECT * FROM accounts;

-- name: name :exec
SELECT accounts.id AS "accounts.id",
     accounts.first_name AS "accounts.first_name",
     accounts.last_name AS "accounts.last_name",
     accounts.email AS "accounts.email",
     accounts.created_at AS "accounts.created_at",
     accounts.updated_at AS "accounts.updated_at"
FROM public.accounts;


