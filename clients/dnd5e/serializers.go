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
		AbilityBonuses:             apiAbilityBonusesToAbilityBonuses(input.AbilityBonuses),
	}
}

func apiAbilityBonusesToAbilityBonuses(input []*apiEntities.AbilityBonus) []*entities.AbilityBonus {
	output := make([]*entities.AbilityBonus, len(input))
	for i, apiAbilityBonus := range input {
		output[i] = apiAbilityBonusToAbilityBonus(apiAbilityBonus)
	}

	return output
}

func apiAbilityBonusToAbilityBonus(input *apiEntities.AbilityBonus) *entities.AbilityBonus {
	if input == nil {
		return nil
	}
	if input.AbilityScore == nil {
		return nil
	}

	return &entities.AbilityBonus{
		Attribute: referenceItemKeyToAttribute(input.AbilityScore.Key),
		Bonus:     input.Bonus,
	}
}

func referenceItemKeyToAttribute(input string) entities.Attribute {
	switch input {
	case "str":
		return entities.AttributeStrength
	case "dex":
		return entities.AttributeDexterity
	case "con":
		return entities.AttributeConstitution
	case "int":
		return entities.AttributeIntelligence
	case "wis":
		return entities.AttributeWisdom
	case "cha":
		return entities.AttributeCharisma
	default:
		log.Fatalf("Unknown attribute %s", input)
		return entities.AttributeNone
	}
}

func apiProficienciesToProficiencies(input []*apiEntities.Proficiency) []*entities.Proficiency {
	output := make([]*entities.Proficiency, len(input))
	for i, apiProficiency := range input {
		output[i] = apiProficiencyToProficiency(apiProficiency)
	}

	return output
}

func apiProficiencyToProficiency(input *apiEntities.Proficiency) *entities.Proficiency {
	return &entities.Proficiency{
		Key:  input.Key,
		Name: input.Name,
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
		Type: typeStringToReferenceType(input.Type),
	}
}

func typeStringToReferenceType(input string) entities.ReferenceType {
	switch input {
	case "equipment":
		return entities.ReferenceTypeEquipment
	case "proficiencies":
		return entities.ReferenceTypeProficiency
	case "languages":
		return entities.ReferenceTypeLanguage
	case "ability-scores":
		return entities.ReferenceTypeAbilityScore
	case "skills":
		return entities.ReferenceTypeSkill
	default:
		log.Println("Unknown reference type: ", input)
		return entities.ReferenceTypeUnset
	}
}
