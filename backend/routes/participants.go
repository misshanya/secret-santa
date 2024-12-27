package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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
	userID := int64(r.Context().Value("user_id").(int))

	_, err = a.queries.GetParticipantByUserID(r.Context(), db.GetParticipantByUserIDParams{
		UserID: userID,
		RoomID: int64(roomID),
	})
	if err == nil {
		http.Error(w, "You are already in this room", http.StatusConflict)
		return
	} else if err != pgx.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.queries.CreateParticipant(r.Context(), db.CreateParticipantParams{
		UserID: userID,
		RoomID: int64(roomID),
	})

	if err != nil {
		http.Error(w, "Failed to enter room", http.StatusInternalServerError)
		return
	}

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

func (a *ParticipantsAPI) DistributeParticipants(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := int64(r.Context().Value("user_id").(int))

	room, err := a.queries.GetRoomByID(r.Context(), int64(roomID))
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if userID != room.OwnerID {
		http.Error(w, "You are not allowed to do this", http.StatusForbidden)
		return
	}

	participantsIDs, err := a.queries.GetAllParticipants(r.Context(), int64(roomID))
	if err != nil {
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	if len(participantsIDs) < 2 {
		http.Error(w, "Not enough participants to distribute", http.StatusBadRequest)
		return
	}

	rand.Shuffle(len(participantsIDs), func(i, j int) {
		participantsIDs[i], participantsIDs[j] = participantsIDs[j], participantsIDs[i]
	})

	for i, ID := range participantsIDs {
		givesTo := participantsIDs[(i+1)%len(participantsIDs)]

		err := a.queries.SetGivesTo(r.Context(), db.SetGivesToParams{
			GivesTo: pgtype.Int8{Int64: givesTo, Valid: true},
			ID:      ID,
		})

		if err != nil {
			http.Error(w, "Failed to distribute participants", http.StatusInternalServerError)
			return
		}
	}
}

type ParticipantResponse struct {
	ID      int64       `json:"id"`
	UserID  int64       `json:"user_id"`
	RoomID  int64       `json:"room_id"`
	Wish    pgtype.Text `json:"wish"`
	GivesTo pgtype.Int8 `json:"gives_to"`
}

func (a *ParticipantsAPI) GetParticipants(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := int64(r.Context().Value("user_id").(int))

	room, err := a.queries.GetRoomByID(r.Context(), int64(roomID))
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if userID != room.OwnerID {
		http.Error(w, "You are not allowed to do this", http.StatusForbidden)
		return
	}

	participantsIDs, err := a.queries.GetAllParticipants(r.Context(), int64(roomID))
	if err != nil {
		http.Error(w, "Failed to get participants", http.StatusInternalServerError)
		return
	}

	var response struct {
		Participants []ParticipantResponse `json:"participants"`
	}

	response.Participants = make([]ParticipantResponse, len(participantsIDs))

	for i, ID := range participantsIDs {
		participant, err := a.queries.GetParticipantByID(r.Context(), ID)
		if err != nil {
			http.Error(w, "Failed to get participant", http.StatusInternalServerError)
			return
		}
		response.Participants[i] = ParticipantResponse{
			ID:      participant.ID,
			UserID:  participant.UserID,
			RoomID:  participant.RoomID,
			Wish:    participant.Wish,
			GivesTo: participant.GivesTo,
		}
	}

	json.NewEncoder(w).Encode(response)
}
