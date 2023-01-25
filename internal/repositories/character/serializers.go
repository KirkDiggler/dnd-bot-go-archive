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
		Str: 0,
		Dex: 0,
		Con: 0,
		Int: 0,
		Wis: 0,
		Cha: 0,
	}

	for key, attr := range input.Attribues {
		switch key {
		case entities.AttributeStrength:
			data.Str = attr.Score
		case entities.AttributeDexterity:
			data.Dex = attr.Score
		case entities.AttributeConstitution:
			data.Con = attr.Score
		case entities.AttributeIntelligence:
			data.Int = attr.Score
		case entities.AttributeWisdom:
			data.Wis = attr.Score
		case entities.AttributeCharisma:
			data.Cha = attr.Score
		}
	}

	proficiencies := make([]*Proficiency, len(input.Proficiencies))
	for i, prof := range input.Proficiencies {
		proficiencies[i] = proficiencyToData(prof)
	}

	return &Data{
		ID:            input.ID,
		OwnerID:       input.OwnerID,
		Name:          input.Name,
		RaceKey:       raceKey,
		ClassKey:      classKey,
		Attributes:    data,
		Rolls:         rollResultsToRollDatas(input.Rolls),
		Proficiencies: proficiencies,
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

func proficiencyToData(input *entities.Proficiency) *Proficiency {
	return &Proficiency{
		Key: input.Key,
	}
}
