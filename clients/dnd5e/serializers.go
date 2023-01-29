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
		Key:                        input.Key,
		Name:                       input.Name,
		StartingProficiencyOptions: apiChoiceOptionToChoice(input.StartingProficiencyOptions),
	}
}

func apiClassToClass(input *apiEntities.Class) *entities.Class {
	return &entities.Class{
		Key:                      input.Key,
		Name:                     input.Name,
		ProficiencyChoices:       apiChoicesToChoices(input.ProficiencyChoices),
		StartingEquipmentChoices: apiChoicesToChoices(input.StartingEquipmentOptions),
	}
}

func apiChoicesToChoices(input []*apiEntities.ChoiceOption) []*entities.Choice {
	output := make([]*entities.Choice, len(input))
	for i, apiChoice := range input {
		output[i] = apiChoiceOptionToChoice(apiChoice)
	}

	return output
}

func apiChoiceOptionToChoice(input *apiEntities.ChoiceOption) *entities.Choice {
	if input == nil {
		return nil
	}

	if input.OptionList == nil {
		return nil
	}

	output := make([]entities.Option, len(input.OptionList.Options))

	for i, apiProficiency := range input.OptionList.Options {
		output[i] = apiOptionToOption(apiProficiency)
	}

	return &entities.Choice{
		Count:   input.ChoiceCount,
		Name:    input.Description,
		Type:    apiChoiceTypeToChoiceType(input.ChoiceType),
		Key:     "choice",
		Options: output,
	}
}

func apiChoiceTypeToChoiceType(input string) entities.ChoiceType {
	switch input {
	case "proficiencies":
		return entities.ChoiceTypeProficiency
	case "equipment":
		return entities.ChoiceTypeEquipment
	case "languages":
		return entities.ChoiceTypeLanguage
	default:
		log.Println("Unknown choice type: ", input)
		return entities.ChoiceTypeUnset
	}
}

func apiOptionToOption(input apiEntities.Option) entities.Option {
	switch input.GetOptionType() {
	case apiEntities.OptionTypeReference:
		item := input.(*apiEntities.ReferenceOption)
		if item.Reference == nil {
			return nil
		}

		return &entities.ReferenceOption{
			Reference: apiReferenceItemToReferenceItem(item.Reference),
		}
	case apiEntities.OptionTypeChoice:
		item := input.(*apiEntities.ChoiceOption)

		return apiChoiceOptionToChoice(item)
	case apiEntities.OptionalTypeCountedReference:
		item := input.(*apiEntities.CountedReferenceOption)
		if item.Reference == nil {
			return nil
		}

		return &entities.CountedReferenceOption{
			Count:     item.Count,
			Reference: apiReferenceItemToReferenceItem(item.Reference),
		}
	case apiEntities.OptionTypeMultiple:
		item := input.(*apiEntities.MultipleOption)
		if item.Items == nil {
			return nil
		}

		options := make([]entities.Option, len(item.Items))
		for i, apiOption := range item.Items {
			options[i] = apiOptionToOption(apiOption)
		}

		return &entities.MultipleOption{
			Items: options,
		}
	default:
		log.Println("Unknown option type: ", input.GetOptionType())
		return nil
	}
}

func apiReferenceItemToReferenceItem(input *apiEntities.ReferenceItem) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Key,
		Name: input.Name,
		Type: input.Type,
	}
}
