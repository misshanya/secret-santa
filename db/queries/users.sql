-- name: RegisterUser :exec
INSERT INTO users (
    name,
    username,
    password
) VALUES (
    $1, $2, $3
);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;