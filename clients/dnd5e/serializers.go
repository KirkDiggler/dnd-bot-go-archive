package dnd5e

import (
	"log"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	apiEntities "github.com/fadedpez/dnd5e-api/entities"
)

func apiReferenceItemToClass(apiClass *apiEntities.ReferenceItem) *entities.Class {
	return &entities.Class{
		Key:  apiClass.Key,
		Name: apiClass.Name,
	}
}

func apiReferenceItemsToClasses(input []*apiEntities.ReferenceItem) []*entities.Class {
	output := make([]*entities.Class, len(input))
	for i, apiClass := range input {
		output[i] = apiReferenceItemToClass(apiClass)
	}
	return output
}

func apiReferenceItemToRace(input *apiEntities.ReferenceItem) *entities.Race {
	return &entities.Race{
		Key:  input.Key,
		Name: input.Name,
	}
}

func apiReferenceItemsToRaces(input []*apiEntities.ReferenceItem) []*entities.Race {
	output := make([]*entities.Race, len(input))
	for i, apiRace := range input {
		output[i] = apiReferenceItemToRace(apiRace)
	}

	return output
}

func apiRaceToRace(input *apiEntities.Race) *entities.Race {
	return &entities.Race{
		Key:  input.Key,
		Name: input.Name,
	}
}

func apiClassToClass(input *apiEntities.Class) *entities.Class {
	return &entities.Class{
		Key:                input.Key,
		Name:               input.Name,
		ProficiencyChoices: apiChoicesToProficiencyChoices(input.ProficiencyChoices),
	}
}

func apiChoicesToProficiencyChoices(input []*apiEntities.ChoiceOption) []*entities.ProficiencyChoices {
	output := make([]*entities.ProficiencyChoices, len(input))
	for i, apiChoice := range input {
		output[i] = apiChoiceOptionsToProcicincies(apiChoice)
	}

	return output
}

func apiChoiceOptionsToProcicincies(input *apiEntities.ChoiceOption) *entities.ProficiencyChoices {
	if input == nil {
		return nil
	}

	if input.OptionList == nil {
		return nil
	}

	output := make([]*entities.Proficiency, len(input.OptionList.Options))
	for i, apiProficiency := range input.OptionList.Options {
		output[i] = apiChoiceOptionToProficiency(apiProficiency)
	}

	return &entities.ProficiencyChoices{
		Count: input.ChoiceCount,
		From:  output,
	}
}

func apiChoiceOptionToProficiency(input apiEntities.Option) *entities.Proficiency {
	switch input.GetOptionType() {
	case apiEntities.OptionTypeReference:
		item := input.(*apiEntities.ReferenceOption)
		if item.Reference == nil {
			return nil
		}

		return &entities.Proficiency{
			Key:  item.Reference.Key,
			Name: item.Reference.Name,
		}
	default:
		log.Println("Unknown option type: ", input.GetOptionType())
		return nil
	}
}
