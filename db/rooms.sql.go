// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rooms.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRoom = `-- name: CreateRoom :exec
INSERT INTO rooms (
    owner_id,
    name,
    description
) VALUES (
    $1,
    $2,
    $3
)
`

type CreateRoomParams struct {
	OwnerID     int64
	Name        pgtype.Text
	Description pgtype.Text
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) error {
	_, err := q.db.Exec(ctx, createRoom, arg.OwnerID, arg.Name, arg.Description)
	return err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteRoom, id)
	return err
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, owner_id, name, description, max_participants, created_at FROM rooms
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRoomByID(ctx context.Context, id int64) (Room, error) {
	row := q.db.QueryRow(ctx, getRoomByID, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Name,
		&i.Description,
		&i.MaxParticipants,
		&i.CreatedAt,
	)
	return i, err
}

const getUserRooms = `-- name: GetUserRooms :many
SELECT id, name, description, max_participants, created_at 
FROM rooms
WHERE owner_id = $1
ORDER BY created_at DESC
`

type GetUserRoomsRow struct {
	ID              int64
	Name            pgtype.Text
	Description     pgtype.Text
	MaxParticipants pgtype.Int4
	CreatedAt       pgtype.Timestamp
}

func (q *Queries) GetUserRooms(ctx context.Context, ownerID int64) ([]GetUserRoomsRow, error) {
	rows, err := q.db.Query(ctx, getUserRooms, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserRoomsRow
	for rows.Next() {
		var i GetUserRoomsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MaxParticipants,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
