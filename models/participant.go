package models

type Participant struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	RoomID uint   `json:"room_id"`
	Wish   string `json:"wish"`
}
