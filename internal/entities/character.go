package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

type Character struct {
	ID      string
	OwnerID string
	Name    string
	Race    *Race
	Class   *Class
}

func (c *Character) ToData() *character.Data {
	var raceKey string
	var classKey string

	if c.Race != nil {
		raceKey = c.Race.Key
	}

	if c.Class != nil {
		classKey = c.Class.Key
	}

	return &character.Data{
		ID:       c.ID,
		OwnerID:  c.OwnerID,
		Name:     c.Name,
		RaceKey:  raceKey,
		ClassKey: classKey,
	}
}
