package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

type Character struct {
	ID        string
	OwnerID   string
	Name      string
	Race      *Race
	Class     *Class
	Attribues map[Attribute]*AbilityScore
	Rolls     []*dice.RollResult
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

	data := &character.AttributeData{
		Str: 0,
		Dex: 0,
		Con: 0,
		Int: 0,
		Wis: 0,
		Cha: 0,
	}

	for key, attr := range c.Attribues {
		switch key {
		case AttributeStrength:
			data.Str = attr.Score
		case AttributeDexterity:
			data.Dex = attr.Score
		case AttributeConstitution:
			data.Con = attr.Score
		case AttributeIntelligence:
			data.Int = attr.Score
		case AttributeWisdom:
			data.Wis = attr.Score
		case AttributeCharisma:
			data.Cha = attr.Score
		}
	}
	return &character.Data{
		ID:         c.ID,
		OwnerID:    c.OwnerID,
		Name:       c.Name,
		RaceKey:    raceKey,
		ClassKey:   classKey,
		Attributes: data,
		Rolls:      rollResultsToRollDatas(c.Rolls),
	}
}

func rollResultToRollData(result *dice.RollResult) *character.RollData {
	if result == nil {

		return nil
	}
	return &character.RollData{
		Used:    result.Used,
		Total:   result.Total,
		Highest: result.Highest,
		Lowest:  result.Lowest,
		Rolls:   result.Rolls,
	}
}

func rollResultsToRollDatas(results []*dice.RollResult) []*character.RollData {
	data := make([]*character.RollData, len(results))
	for i, r := range results {
		data[i] = rollResultToRollData(r)
	}

	return data
}
