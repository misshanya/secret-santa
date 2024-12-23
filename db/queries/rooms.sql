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
SELECT id, name, description, max_participants, created_at 
FROM rooms
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;