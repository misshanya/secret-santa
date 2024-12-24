package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/misshanya/secret-santa/db"
)

type ParticipantsAPI struct {
	queries *db.Queries
}

func NewParticipantsAPI(queries *db.Queries) *ParticipantsAPI {
	return &ParticipantsAPI{queries: queries}
}

func (a *ParticipantsAPI) NewParticipant(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	UserID := int64(r.Context().Value("user_id").(int))

	a.queries.CreateParticipant(r.Context(), db.CreateParticipantParams{
		UserID: UserID,
		RoomID: int64(roomID),
	})

	w.WriteHeader(http.StatusCreated)
}

func (a *ParticipantsAPI) GetWish(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	userID := int64(r.Context().Value("user_id").(int))

	participant, err := a.queries.GetParticipantByUserID(r.Context(), db.GetParticipantByUserIDParams{
		UserID: userID,
		RoomID: int64(roomID),
	})
	if err != nil {
		http.Error(w, "You are not member of this room", http.StatusNotFound)
		return
	}

	response := struct {
		Wish string `json:"wish"`
	}{
		Wish: participant.Wish.String,
	}

	json.NewEncoder(w).Encode(response)
}

func (a *ParticipantsAPI) SetWish(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	userID := int64(r.Context().Value("user_id").(int))

	var body struct {
		Wish string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = a.queries.GetParticipantByUserID(r.Context(), db.GetParticipantByUserIDParams{
		UserID: userID,
		RoomID: int64(roomID),
	})
	if err != nil {
		http.Error(w, "You are not a member of this room", http.StatusNotFound)
		return
	}

	a.queries.UpdateParticipiantWish(r.Context(), db.UpdateParticipiantWishParams{
		Wish:   pgtype.Text{String: body.Wish, Valid: true},
		UserID: userID,
		RoomID: int64(roomID),
	})
}

func (a *ParticipantsAPI) DeleteParticipant(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	participant, err := a.queries.GetParticipantByUserID(r.Context(), db.GetParticipantByUserIDParams{
		UserID: int64(r.Context().Value("user_id").(int)),
		RoomID: int64(roomID),
	})
	if err != nil {
		http.Error(w, "This participant does not exists", http.StatusNotFound)
		return
	}

	a.queries.DeleteParticipant(r.Context(), participant.ID)
}
