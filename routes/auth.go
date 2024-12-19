package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (a *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := a.queries.GetUserByUsername(r.Context(), body.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(10 * time.Minute).Unix(),
		})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(tokenString)
}

func (a *AuthAPI) GetMyInfo(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(int))

	user, err := a.queries.GetUserByID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	}{
		ID:       user.ID,
		Name:     user.Name.String,
		Username: user.Username,
	}

	json.NewEncoder(w).Encode(response)
}
