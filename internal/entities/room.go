package entities

import "fmt"

type RoomStatus string

const (
	RoomStatusUnset    RoomStatus = ""
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
)

type Room struct {
	ID        string
	Status    RoomStatus
	Character *Character
	Monster   *Monster
}

func (r *Room) String() string {
	return fmt.Sprintf("room: %s, status: %s, character: %s, monster: %s", r.ID, r.Status, r.Character, r.Monster)
}
