package entities

type RoomStatus string

const (
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
)

type Room struct {
	ID          string     `json:"id"`
	Status      RoomStatus `json:"status"`
	CharacterID string     `json:"character_id"`
	MonsterID   string     `json:"monster_id"`
}
