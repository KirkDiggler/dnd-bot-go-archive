package entities

import (
	"fmt"
)

type RoomStatus string

const (
	RoomStatusUnset    RoomStatus = ""
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
	RoomStatusWon      RoomStatus = "won"
	RoomStatusLost     RoomStatus = "lost"
)

type Room struct {
	ID                  string
	Status              RoomStatus
	Character           *Character
	CharacterInitiative int
	Monster             *Monster
	MonsterInitiative   int
}

func (r *Room) IsEmpty() bool {
	return r.Character == nil && r.Monster == nil
}

func (r *Room) IsFull() bool {
	return r.Character != nil && r.Monster != nil
}

func (r *Room) IsCharacterAlive() bool {
	return r.Character != nil && r.Character.CurrentHitPoints > 0
}

func (r *Room) IsMonsterAlive() bool {
	return r.Monster != nil && r.Monster.CurrentHP > 0
}

func (r *Room) String() string {
	return fmt.Sprintf("room: %s, status: %s\ncharacter: %s\nmonster: %s", r.ID, r.Status, r.Character.ShortString(), r.Monster)
}
