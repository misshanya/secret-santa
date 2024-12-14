package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/misshanya/secret-santa/db"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	if _, err := a.queries.GetUserByUsername(r.Context(), body.Username); err == nil {
		http.Error(w, "This user already exists", http.StatusConflict)
		return
	}

	err = a.queries.RegisterUser(r.Context(), db.RegisterUserParams{
		Name:     pgtype.Text{String: body.Name, Valid: true},
		Username: body.Username,
		Password: string(hashedPassword),
	})

	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
