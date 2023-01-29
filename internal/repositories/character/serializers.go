package character

import (
	"encoding/json"
	"log"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

func jsonToData(input string) *Data {
	char := &Data{}

	err := json.Unmarshal([]byte(input), char)
	if err != nil {
		log.Println(err)
		return nil
	}

	return char
}

func dataToJSON(input *Data) string {
	b, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(b)
}

func characterToData(input *entities.Character) *Data {
	var raceKey string
	var classKey string

	if input.Race != nil {
		raceKey = input.Race.Key
	}

	if input.Class != nil {
		classKey = input.Class.Key
	}

	data := &AttributeData{
		Str: &AbilityScoreData{},
		Dex: &AbilityScoreData{},
		Con: &AbilityScoreData{},
		Int: &AbilityScoreData{},
		Wis: &AbilityScoreData{},
		Cha: &AbilityScoreData{},
	}

	for key, attr := range input.Attribues {
		switch key {
		case entities.AttributeStrength:
			data.Str = abilityScoreToData(attr)
		case entities.AttributeDexterity:
			data.Dex = abilityScoreToData(attr)
		case entities.AttributeConstitution:
			data.Con = abilityScoreToData(attr)
		case entities.AttributeIntelligence:
			data.Int = abilityScoreToData(attr)
		case entities.AttributeWisdom:
			data.Wis = abilityScoreToData(attr)
		case entities.AttributeCharisma:
			data.Cha = abilityScoreToData(attr)
		}
	}

	return &Data{
		ID:            input.ID,
		OwnerID:       input.OwnerID,
		Name:          input.Name,
		RaceKey:       raceKey,
		ClassKey:      classKey,
		Attributes:    data,
		Rolls:         rollResultsToRollDatas(input.Rolls),
		Proficiencies: proficienciesToDatas(input.Proficiencies),
	}
}
func abilityScoreToData(input *entities.AbilityScore) *AbilityScoreData {
	return &AbilityScoreData{
		Score: input.Score,
		Bonus: input.Bonus,
	}
}
func rollResultToRollData(result *dice.RollResult) *RollData {
	if result == nil {

		return nil
	}
	return &RollData{
		Used:    result.Used,
		Total:   result.Total,
		Highest: result.Highest,
		Lowest:  result.Lowest,
		Rolls:   result.Rolls,
	}
}
func rollResultsToRollDatas(results []*dice.RollResult) []*RollData {
	data := make([]*RollData, len(results))
	for i, r := range results {
		data[i] = rollResultToRollData(r)
	}

	return data
}

func proficienciesToDatas(input []*entities.Proficiency) []*Proficiency {
	data := make([]*Proficiency, len(input))
	for i, p := range input {
		data[i] = proficiencyToData(p)
	}

	return data
}

func proficiencyToData(input *entities.Proficiency) *Proficiency {
	return &Proficiency{
		Key:  input.Key,
		Name: input.Name,
	}
}
