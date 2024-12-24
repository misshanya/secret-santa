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

-- name: GetParticipantByUserID :one
SELECT * FROM participants
WHERE user_id = $1 AND room_id = $2;

-- name: GetUserParticipations :many
SELECT * FROM participants
WHERE user_id = $1;

-- name: UpdateParticipiantWish :exec
UPDATE participants
SET wish = $1
WHERE user_id = $2 AND room_id = $3;

-- name: DeleteParticipant :exec
DELETE FROM participants
WHERE id = $1;