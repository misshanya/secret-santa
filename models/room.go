package models

type Room struct {
	ID              uint   `json:"id"`
	OwnerID         uint   `json:"owner_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	MaxParticipants uint   `json:"max_participants"`
}
