-- name: CreateUser :one
INSERT INTO users (
    username,
    password,
    email
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET
    username = $2,
    password = $3,
    email = $4
WHERE id = $1;

-- name: GetUser :one
SELECT username,email FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserLoginParams :one
SELECT id, password 
FROM users 
WHERE username = $1;

-- name: UserExistsByID :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);