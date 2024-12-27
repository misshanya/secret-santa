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

-- name: GetAllParticipants :many
SELECT id FROM participants
WHERE room_id = $1
ORDER BY id DESC;

-- name: SetGivesTo :exec
UPDATE participants
SET gives_to = $1
WHERE id = $2;

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