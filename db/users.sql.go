// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, username, password, registered_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Password,
		&i.RegisteredAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, name, username, password, registered_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.Password,
		&i.RegisteredAt,
	)
	return i, err
}

const registerUser = `-- name: RegisterUser :exec
INSERT INTO users (
    name,
    username,
    password
) VALUES (
    $1, $2, $3
)
`

type RegisterUserParams struct {
	Name     pgtype.Text
	Username string
	Password string
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) error {
	_, err := q.db.Exec(ctx, registerUser, arg.Name, arg.Username, arg.Password)
	return err
}
