package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
		OwnerID:     int64(r.Context().Value("user_id").(int)),
		Name:        pgtype.Text{String: body.Name, Valid: true},
		Description: pgtype.Text{String: body.Description, Valid: true},
	})

	w.WriteHeader(http.StatusCreated)
}

type RoomResponse struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	MaxParticipants int32  `json:"max_participants"`
	CreatedAt       string `json:"created_at"`
}

func (a *RoomsAPI) MyRooms(w http.ResponseWriter, r *http.Request) {
	userID := int64(r.Context().Value("user_id").(int))

	rooms, err := a.queries.GetUserRooms(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Rooms []RoomResponse `json:"rooms"`
	}

	for _, room := range rooms {
		response.Rooms = append(response.Rooms, RoomResponse{
			ID:              room.ID,
			Name:            room.Name.String,
			Description:     room.Description.String,
			MaxParticipants: room.MaxParticipants.Int32,
			CreatedAt:       room.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (a *RoomsAPI) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	room, err := a.queries.GetRoomByID(r.Context(), int64(roomID))
	if err != nil {
		http.Error(w, "This room does not exists", http.StatusNotFound)
		return
	}

	if room.OwnerID != int64(r.Context().Value("user_id").(int)) {
		http.Error(w, "You are not allowed to do this", http.StatusForbidden)
		return
	}

	a.queries.DeleteRoom(r.Context(), int64(roomID))
}
