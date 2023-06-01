package encounter

import "time"

type Status string

const (
	StatusUnset    Status = ""
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Data struct {
	ID        string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    Status    `json:"status"`
	PlayerID  string    `json:"player_id"`
	MonsterID string    `json:"monster_id"`
	RoomID    string    `json:"room_id"`
}
