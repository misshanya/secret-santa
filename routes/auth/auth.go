package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/misshanya/secret-santa/db"
)

type AuthAPI struct {
	queries *db.Queries
}

func NewAuthAPI(queries *db.Queries) *AuthAPI {
	return &AuthAPI{queries: queries}
}

func (a *AuthAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	err := a.queries.RegisterUser(ctx, db.RegisterUserParams{
		Name:     pgtype.Text{String: body.Name, Valid: true},
		Username: body.Username,
		Password: body.Password,
	})

	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
