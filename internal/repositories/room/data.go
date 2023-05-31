package room

import "github.com/KirkDiggler/dnd-bot-go/internal/entities"

type RoomStatus string

const (
	RoomStatusUnset    RoomStatus = ""
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
)

type Data struct {
	ID          string     `json:"id"`
	Status      RoomStatus `json:"status"`
	CharacterID string     `json:"character_id"`
	MonsterID   string     `json:"monster_id"`
}

func EntityToRoomStatus(input entities.RoomStatus) RoomStatus {
	switch input {
	case entities.RoomStatusActive:
		return RoomStatusActive
	case entities.RoomStatusInactive:
		return RoomStatusInactive
	default:
		return RoomStatusUnset
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
		ID:          input.ID,
		Status:      EntityToRoomStatus(input.Status),
		CharacterID: charID,
		MonsterID:   monsterID,
	}
}
