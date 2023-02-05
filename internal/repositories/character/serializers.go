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
		ID:               input.ID,
		OwnerID:          input.OwnerID,
		Name:             input.Name,
		HitDie:           input.HitDie,
		AC:               input.AC,
		MaxHitPoints:     input.MaxHitPoints,
		CurrentHitPoints: input.CurrentHitPoints,
		Experience:       input.Experience,
		NextLevel:        input.NextLevel,
		Speed:            input.Speed,
		Level:            input.Level,
		RaceKey:          raceKey,
		ClassKey:         classKey,
		Attributes:       data,
		Rolls:            rollResultsToRollDatas(input.Rolls),
		Proficiencies:    proficienciesToDatas(input.Proficiencies),
		Inventory:        equipmentsToDatas(input.Inventory),
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

func proficienciesToDatas(input map[entities.ProficiencyType][]*entities.Proficiency) []*Proficiency {
	datas := make([]*Proficiency, 0)

	for _, v := range input {
		for _, p := range v {
			datas = append(datas, proficiencyToData(p))
		}
	}

	return datas
}

func proficiencyToData(input *entities.Proficiency) *Proficiency {
	return &Proficiency{
		Key:  input.Key,
		Name: input.Name,
		Type: string(input.Type),
	}
}

func equipmentToData(input entities.Equipment) *Equipment {
	return &Equipment{
		Key:  input.GetKey(),
		Name: input.GetName(),
		Type: input.GetEquipmentType(),
	}
}

func equipmentsToDatas(input map[string][]entities.Equipment) []*Equipment {
	datas := make([]*Equipment, 0)

	for _, v := range input {
		for _, e := range v {
			log.Println("adding equipment: ", e.GetName())
			datas = append(datas, equipmentToData(e))
		}
	}

	return datas
}
