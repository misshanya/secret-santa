-- name: CreateParticipant :exec
INSERT INTO participants (
    user_id,
    room_id,
    wish
) VALUES (
    $1,
    $2,
    $3
);

-- name: GetParticipantByID :one
SELECT * FROM participants
WHERE id = $1;

-- name: GetUserParticipations :many
SELECT * FROM participants
WHERE user_id = $1;

-- name: DeleteParticipant :exec
DELETE FROM participants
WHERE id = $1;