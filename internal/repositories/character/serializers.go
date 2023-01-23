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
		ID:       input.ID,
		OwnerID:  input.OwnerID,
		Name:     input.Name,
		RaceKey:  raceKey,
		ClassKey: classKey,
		Class: &ClassData{
			Key:                classKey,
			ProficiencyChoices: input.Class.ProficiencyChoices,
		},
		Attributes:    data,
		Rolls:         rollResultsToRollDatas(input.Rolls),
		Proficiencies: proficiencies,
	}
}
func multipleOptionToData(input *entities.MultipleOption) *MultipleOption {
	return &MultipleOption{
		Selected: input.Selected,
		Items:    optionsToDatas(input.Items),
	}
}

func optionsToDatas(input []entities.Option) []*Option {
	out := make([]*Option, len(input))
	for i, opt := range input {
		out[i] = optionToData(opt)
	}

	return out
}
func optionToData(input entities.Option) *Option {
	switch input.GetOptionType() {
	case entities.OptionTypeChoice:
		return &Option{
			Choice: ChoiceToData(input.(*entities.Choice)),
		}
	case entities.OptionTypeMultiple:
		return &Option{
			Multiple: multipleOptionToData(input.(*entities.MultipleOption)),
		}
	case entities.OptionTypeReference:
		return &Option{
			Reference: referenceItemToData(input.(*entities.ReferenceOption)),
		}
	case entities.OptionTypeCountedReference:
		return &Option{
			CountedReference: countedReferenceItemToData(input.(*entities.CountedReferenceOption)),
		}
	default:
		log.Println("Unknown option type")
		return nil
	}
}

func countedReferenceItemToData(input *entities.CountedReferenceOption) *CountedReferenceOption {
	if input == nil {
		return nil
	}

	if input.Reference == nil {
		return nil
	}

	return &CountedReferenceOption{
		Count: input.Count,
		Reference: &ReferenceItem{
			Key: input.Reference.Key,
		},
	}
}
func referenceItemToData(input *entities.ReferenceOption) *ReferenceOption {
	if input == nil {
		return nil
	}

	if input.Reference == nil {
		return nil
	}

	return &ReferenceOption{
		&ReferenceItem{
			Key: input.Reference.Key,
		},
	}
}
func ChoiceToData(input *entities.Choice) *Choice {
	choice := &Choice{
		Count:    input.Count,
		Selected: input.Selected,
		Name:     input.Name,
	}

	choice.Options = make([]*Option, len(input.Options))
	for i, opt := range input.Options {
		choice.Options[i] = optionToData(opt)
	}

	return choice
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
