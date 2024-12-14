package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/misshanya/secret-santa/db"
)

type RoomsAPI struct {
	queries *db.Queries
}

func NewRoomsAPI(queries *db.Queries) *RoomsAPI {
	return &RoomsAPI{queries: queries}
}

func (a *RoomsAPI) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string
		Description string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	a.queries.CreateRoom(r.Context(), db.CreateRoomParams{
		OwnerID:     int32(r.Context().Value("user_id").(int)),
		Name:        pgtype.Text{String: body.Name, Valid: true},
		Description: pgtype.Text{String: body.Description, Valid: true},
	})

	w.WriteHeader(http.StatusCreated)
}

func (a *RoomsAPI) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID int
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	room, err := a.queries.GetRoomByID(r.Context(), int32(body.ID))
	if err != nil {
		http.Error(w, "This room does not exists", http.StatusNotFound)
		return
	}

	if room.OwnerID != int32(r.Context().Value("user_id").(int)) {
		http.Error(w, "You are not allowed to do this", http.StatusForbidden)
		return
	}

	a.queries.DeleteRoom(r.Context(), int32(body.ID))
}
