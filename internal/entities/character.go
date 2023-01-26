package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
)

type Character struct {
	ID                 string
	OwnerID            string
	Name               string
	Race               *Race
	Class              *Class
	Attribues          map[Attribute]*AbilityScore
	Rolls              []*dice.RollResult
	Proficiencies      []*Proficiency
	ProficiencyChoices []*Choice
}

func (c *Character) AddProficiency(p *Proficiency) {
	if c.Proficiencies == nil {
		c.Proficiencies = make([]*Proficiency, 0)
	}

	c.Proficiencies = append(c.Proficiencies, p)
}
