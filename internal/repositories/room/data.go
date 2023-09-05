package room

import "github.com/KirkDiggler/dnd-bot-go/internal/entities"

type Status string

const (
	StatusUnset    Status = ""
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusWon      Status = "won"
	StatusLost     Status = "lost"
)

type Data struct {
	ID                string `json:"id"`
	Status            Status `json:"status"`
	PlayerID          string `json:"player_id"`
	PlayerInitiative  int    `json:"player_initiative"`
	MonsterID         string `json:"monster_id"`
	MonsterInitiative int    `json:"monster_initiative"`
}

func EntityToRoomStatus(input entities.RoomStatus) Status {
	switch input {
	case entities.RoomStatusActive:
		return StatusActive
	case entities.RoomStatusInactive:
		return StatusInactive
	case entities.RoomStatusWon:
		return StatusWon
	case entities.RoomStatusLost:
		return StatusLost
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
		ID:                input.ID,
		Status:            EntityToRoomStatus(input.Status),
		PlayerID:          charID,
		PlayerInitiative:  input.CharacterInitiative,
		MonsterID:         monsterID,
		MonsterInitiative: input.MonsterInitiative,
	}
}
