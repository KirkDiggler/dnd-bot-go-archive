package entities

import (
	"fmt"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
)

type Character struct {
	ID                 string
	OwnerID            string
	Name               string
	Speed              int
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

func (c *Character) AddAbilityScoreBonus(attr Attribute, bonus int) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	c.Attribues[attr] = c.Attribues[attr].AddBonus(bonus)
}

func (c *Character) Display() string {
	msg := strings.Builder{}
	if c.Race == nil || c.Class == nil {
		return "Character not fully created"
	}

	msg.WriteString(fmt.Sprintf("Name: %s the %s %s\n", c.Name, c.Race.Name, c.Class.Name))

	msg.WriteString("**Rolls**:\n")
	for _, roll := range c.Rolls {
		msg.WriteString(fmt.Sprintf("%s,", roll.Display()))
	}

	msg.WriteString("\n**Attributes**:\n")
	for attr, score := range c.Attribues {
		msg.WriteString(fmt.Sprintf("%s: %s\n", attr, score.Display()))
	}

	return msg.String()
}
