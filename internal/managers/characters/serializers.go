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
		ID:            data.ID,
		Name:          data.Name,
		OwnerID:       data.OwnerID,
		Race:          race,
		Class:         class,
		Attribues:     attributDataToAttributes(data.Attributes),
		Rolls:         rollDatasToRollResults(data.Rolls),
		Proficiencies: datasToProficiencies(data.Proficiencies),
	}, nil
}

func datasToProficiencies(data []*character.Proficiency) []*entities.Proficiency {
	if data == nil {
		return nil
	}

	proficiencies := make([]*entities.Proficiency, 0, len(data))
	for _, d := range data {
		proficiencies = append(proficiencies, dataToProficiency(d))
	}

	return proficiencies

}

func dataToProficiency(data *character.Proficiency) *entities.Proficiency {
	if data == nil {
		return nil
	}

	return &entities.Proficiency{
		Name: data.Name,
		Key:  data.Key,
	}
}
func attributDataToAttributes(data *character.AttributeData) map[entities.Attribute]*entities.AbilityScore {
	if data == nil {
		data = &character.AttributeData{}
	}

	return map[entities.Attribute]*entities.AbilityScore{
		entities.AttributeStrength:     abilityScoreDataToAbilityScore(data.Str),
		entities.AttributeDexterity:    abilityScoreDataToAbilityScore(data.Dex),
		entities.AttributeConstitution: abilityScoreDataToAbilityScore(data.Con),
		entities.AttributeIntelligence: abilityScoreDataToAbilityScore(data.Int),
		entities.AttributeWisdom:       abilityScoreDataToAbilityScore(data.Wis),
		entities.AttributeCharisma:     abilityScoreDataToAbilityScore(data.Cha),
	}
}

func abilityScoreDataToAbilityScore(data *character.AbilityScoreData) *entities.AbilityScore {
	if data == nil {
		return nil
	}

	return &entities.AbilityScore{
		Score: data.Score,
		Bonus: data.Bonus,
	}
}
func rollDataToRollResult(data *character.RollData) *dice.RollResult {
	if data == nil {
		return nil
	}

	return &dice.RollResult{
		Used:    data.Used,
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

func characterToData(input *entities.Character) *character.Data {
	if input == nil {
		return nil
	}

	return &character.Data{
		ID:         input.ID,
		Name:       input.Name,
		OwnerID:    input.OwnerID,
		RaceKey:    input.Race.Key,
		ClassKey:   input.Class.Key,
		Attributes: attributesToAttributeData(input.Attribues),
		Rolls:      rollResultsToRollDatas(input.Rolls),
	}
}

func rollResultsToRollDatas(input []*dice.RollResult) []*character.RollData {
	datas := make([]*character.RollData, len(input))
	for i, r := range input {
		datas[i] = rollResultToRollData(r)
	}

	return datas
}
func rollResultToRollData(input *dice.RollResult) *character.RollData {
	if input == nil {
		return nil
	}

	return &character.RollData{
		Used:    input.Used,
		Total:   input.Total,
		Highest: input.Highest,
		Lowest:  input.Lowest,
		Rolls:   input.Rolls,
	}
}

func abilityScoreToData(input *entities.AbilityScore) *character.AbilityScoreData {
	if input == nil {
		return nil
	}

	return &character.AbilityScoreData{
		Score: input.Score,
		Bonus: input.Bonus,
	}
}
func attributesToAttributeData(input map[entities.Attribute]*entities.AbilityScore) *character.AttributeData {
	if input == nil {
		return nil
	}

	return &character.AttributeData{
		Str: abilityScoreToData(input[entities.AttributeStrength]),
		Dex: abilityScoreToData(input[entities.AttributeDexterity]),
		Con: abilityScoreToData(input[entities.AttributeConstitution]),
		Int: abilityScoreToData(input[entities.AttributeIntelligence]),
		Wis: abilityScoreToData(input[entities.AttributeWisdom]),
		Cha: abilityScoreToData(input[entities.AttributeCharisma]),
	}
}
