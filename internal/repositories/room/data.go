package room

import "github.com/KirkDiggler/dnd-bot-go/internal/entities"

type Status string

const (
	StatusUnset    Status = ""
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Data struct {
	ID        string `json:"id"`
	Status    Status `json:"status"`
	PlayerID  string `json:"player_id"`
	MonsterID string `json:"monster_id"`
}

func EntityToRoomStatus(input entities.RoomStatus) Status {
	switch input {
	case entities.RoomStatusActive:
		return StatusActive
	case entities.RoomStatusInactive:
		return StatusInactive
	default:
		return StatusUnset
	}
}

func EntityToData(input *entities.Room) *Data {
	if input == nil {
		return nil
	}

	var charID, monsterID string
	if input.Character != nil {
		charID = input.Character.ID
	}

	if input.Monster != nil {
		monsterID = input.Monster.ID
	}

	return &Data{
		ID:        input.ID,
		Status:    EntityToRoomStatus(input.Status),
		PlayerID:  charID,
		MonsterID: monsterID,
	}
}
