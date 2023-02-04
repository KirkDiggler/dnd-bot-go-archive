package entities

import (
	"fmt"
	"strings"
	"sync"

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
	Proficiencies      map[ProficiencyType][]*Proficiency
	ProficiencyChoices []*Choice
	mu                 sync.Mutex
}

func (c *Character) AddAttribute(attr Attribute, score int) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	bonus := 0
	if _, ok := c.Attribues[attr]; ok {
		bonus = c.Attribues[attr].Bonus
	}
	abilityScore := &AbilityScore{
		Score: score,
		Bonus: bonus,
	}
	switch {
	case score == 1:
		abilityScore.Bonus += -5
	case score < 4 && score > 1:
		abilityScore.Bonus += -4
	case score < 6 && score > 3:
		abilityScore.Bonus += -3
	case score < 8 && score > 5:
		abilityScore.Bonus += -2
	case score < 10 && score >= 8:
		abilityScore.Bonus += -1
	case score < 12 && score > 9:
		abilityScore.Bonus += 0
	case score < 14 && score > 11:
		abilityScore.Bonus += 1
	case score < 16 && score > 13:
		abilityScore.Bonus += 2
	case score < 18 && score > 15:
		abilityScore.Bonus += 3
	case score < 20 && score > 17:
		abilityScore.Bonus += 4
	case score == 20:
		abilityScore.Bonus += 5
	}

	c.Attribues[attr] = abilityScore
}
func (c *Character) AddAbilityBonus(ab *AbilityBonus) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	if _, ok := c.Attribues[ab.Attribute]; !ok {
		c.Attribues[ab.Attribute] = &AbilityScore{}
	}

	c.Attribues[ab.Attribute] = c.Attribues[ab.Attribute].AddBonus(ab.Bonus)
}

func (c *Character) AddProficiency(p *Proficiency) {
	if c.Proficiencies == nil {
		c.Proficiencies = make(map[ProficiencyType][]*Proficiency)
	}
	c.mu.Lock()
	if c.Proficiencies[p.Type] == nil {
		c.Proficiencies[p.Type] = make([]*Proficiency, 0)
	}

	c.Proficiencies[p.Type] = append(c.Proficiencies[p.Type], p)
	c.mu.Unlock()
}

func (c *Character) AddAbilityScoreBonus(attr Attribute, bonus int) {
	if c.Attribues == nil {
		c.Attribues = make(map[Attribute]*AbilityScore)
	}

	c.Attribues[attr] = c.Attribues[attr].AddBonus(bonus)
}

func (c *Character) String() string {
	msg := strings.Builder{}
	if c.Race == nil || c.Class == nil {
		return "Character not fully created"
	}

	msg.WriteString(fmt.Sprintf("%s the %s %s\n", c.Name, c.Race.Name, c.Class.Name))

	msg.WriteString("**Rolls**:\n")
	for _, roll := range c.Rolls {
		msg.WriteString(fmt.Sprintf("%s, ", roll))
	}

	msg.WriteString("\n**Attributes**:\n")
	for _, attr := range Attributes {
		if c.Attribues[attr] == nil {
			continue
		}
		msg.WriteString(fmt.Sprintf("  -  %s: %s\n", attr, c.Attribues[attr]))
	}

	msg.WriteString("\n**Proficiencies**:\n")
	for _, key := range ProficiencyTypes {
		if c.Proficiencies[key] == nil {
			continue
		}

		msg.WriteString(fmt.Sprintf("  -  **%s**:\n", key))
		for _, prof := range c.Proficiencies[key] {
			msg.WriteString(fmt.Sprintf("    -  %s\n", prof.Name))
		}
	}

	return msg.String()
}
