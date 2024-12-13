-- name: CreateRoom :exec
INSERT INTO rooms (
    owner_id,
    name,
    description
) VALUES (
    $1,
    $2,
    $3
);

-- name: GetRoomByID :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;

-- name: GetUserRooms :many
SELECT * FROM rooms
WHERE owner_id = $1;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;