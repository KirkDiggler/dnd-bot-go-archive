package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

func (m *manager) characterFromData(ctx context.Context, data *character.Data) (*entities.Character, error) {
	if data == nil {
		return nil, dnderr.NewMissingParameterError("data")
	}

	race, err := m.client.GetRace(data.RaceKey)
	if err != nil {
		return nil, err
	}

	class, err := m.client.GetClass(data.ClassKey)
	if err != nil {
		return nil, err
	}

	return &entities.Character{
		ID:        data.ID,
		Name:      data.Name,
		OwnerID:   data.OwnerID,
		Race:      race,
		Class:     class,
		Attribues: attributDataToAttributes(data.Attributes),
		Rolls:     rollDatasToRollResults(data.Rolls),
	}, nil
}

func attributDataToAttributes(data *character.AttributeData) map[entities.Attribute]*entities.AbilityScore {
	if data == nil {
		data = &character.AttributeData{}
	}

	return map[entities.Attribute]*entities.AbilityScore{
		entities.AttributeStrength:     {Score: data.Str},
		entities.AttributeDexterity:    {Score: data.Dex},
		entities.AttributeConstitution: {Score: data.Con},
		entities.AttributeIntelligence: {Score: data.Int},
		entities.AttributeWisdom:       {Score: data.Wis},
		entities.AttributeCharisma:     {Score: data.Cha},
	}
}

func rollDataToRollResult(data *character.RollData) *dice.RollResult {
	if data == nil {
		return nil
	}

	return &dice.RollResult{
		Total:   data.Total,
		Highest: data.Highest,
		Lowest:  data.Lowest,
		Rolls:   data.Rolls,
	}
}

func rollDatasToRollResults(data []*character.RollData) []*dice.RollResult {
	results := make([]*dice.RollResult, len(data))
	for i, d := range data {
		results[i] = rollDataToRollResult(d)
	}

	return results
}
