package choice

import (
	"encoding/json"
	"log"

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

func multipleOptionToData(input *entities.MultipleOption) *MultipleOption {
	return &MultipleOption{
		Status:  choiceStatusToStatus(input.Status),
		Options: optionsToDatas(input.Items),
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
			Choice: choiceToData(input.(*entities.Choice)),
		}
	case entities.OptionTypeMultiple:
		return &Option{
			Multiple: multipleOptionToData(input.(*entities.MultipleOption)),
		}
	case entities.OptionTypeReference:
		return &Option{
			Reference: referenceOptionToData(input.(*entities.ReferenceOption)),
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
		Status:    choiceStatusToStatus(input.Status),
		Count:     input.Count,
		Reference: referenceItemToData(input.Reference),
	}
}
func referenceOptionToData(input *entities.ReferenceOption) *ReferenceOption {
	if input == nil {
		return nil
	}

	if input.Reference == nil {
		return nil
	}

	return &ReferenceOption{
		Status:    choiceStatusToStatus(input.Status),
		Reference: referenceItemToData(input.Reference),
	}
}

func dataToChoice(input *Choice) *entities.Choice {
	choice := &entities.Choice{
		Type:   typeToChoiceType(input.Type),
		Count:  input.Count,
		Status: statusToChoiceStatus(input.Status),
		Name:   input.Name,
	}

	choice.Options = make([]entities.Option, len(input.Options))
	for i, opt := range input.Options {
		choice.Options[i] = dataToOption(opt)
	}

	return choice
}
func choicesToDatas(input []*entities.Choice) []*Choice {
	out := make([]*Choice, len(input))
	for i, choice := range input {
		out[i] = choiceToData(choice)
	}

	return out
}

func choiceToData(input *entities.Choice) *Choice {
	choice := &Choice{
		Type:   choiceTypeToType(input.Type),
		Count:  input.Count,
		Status: choiceStatusToStatus(input.Status),
		Name:   input.Name,
	}

	choice.Options = make([]*Option, len(input.Options))
	for i, opt := range input.Options {
		choice.Options[i] = optionToData(opt)
	}

	return choice
}

func dataToMultipleOption(input *MultipleOption) *entities.MultipleOption {
	return &entities.MultipleOption{
		Status: statusToChoiceStatus(input.Status),
		Items:  datasToOptions(input.Options),
	}
}
func datasToOptions(input []*Option) []entities.Option {
	out := make([]entities.Option, len(input))
	for i, opt := range input {
		out[i] = dataToOption(opt)
	}

	return out
}

func dataToReferenceOption(input *ReferenceOption) *entities.ReferenceOption {
	if input == nil {
		return nil
	}

	if input.Reference == nil {
		return nil
	}

	return &entities.ReferenceOption{
		Status: statusToChoiceStatus(input.Status),
		Reference: &entities.ReferenceItem{
			Key:  input.Reference.Key,
			Name: input.Reference.Name,
			Type: entities.ReferenceType(input.Reference.Type),
		},
	}
}
func dataToCountedReferenceOption(input *CountedReferenceOption) *entities.CountedReferenceOption {
	if input == nil {
		return nil
	}

	if input.Reference == nil {
		return nil
	}

	return &entities.CountedReferenceOption{
		Status:    statusToChoiceStatus(input.Status),
		Count:     input.Count,
		Reference: dataToReferenceItem(input.Reference),
	}
}
func dataToOption(input *Option) entities.Option {
	if input.Choice != nil {
		return dataToChoice(input.Choice)
	}

	if input.Multiple != nil {
		return dataToMultipleOption(input.Multiple)
	}

	if input.Reference != nil {
		return dataToReferenceOption(input.Reference)
	}

	if input.CountedReference != nil {
		return dataToCountedReferenceOption(input.CountedReference)
	}

	return nil
}

func referenceItemToData(input *entities.ReferenceItem) *ReferenceItem {
	if input == nil {
		return nil
	}

	return &ReferenceItem{
		Key:  input.Key,
		Name: input.Name,
		Type: string(input.Type),
	}
}

func dataToReferenceItem(input *ReferenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Key,
		Name: input.Name,
		Type: entities.ReferenceType(input.Type),
	}
}

func datasToChoices(input []*Choice) []*entities.Choice {
	out := make([]*entities.Choice, len(input))
	for i, choice := range input {
		out[i] = dataToChoice(choice)
	}

	return out
}

func choiceStatusToStatus(input entities.ChoiceStatus) Status {
	switch input {
	case entities.ChoiceStatusInactive:
		return StatusInactive
	case entities.ChoiceStatusActive:
		return StatusActive
	case entities.ChoiceStatusSelected:
		return StatusSelected
	default:
		return StatusUnset
	}
}

func statusToChoiceStatus(input Status) entities.ChoiceStatus {
	switch input {
	case StatusInactive:
		return entities.ChoiceStatusInactive
	case StatusActive:
		return entities.ChoiceStatusActive
	case StatusSelected:
		return entities.ChoiceStatusSelected
	default:
		return entities.ChoiceStatusUnset
	}
}

func choiceTypeToType(input entities.ChoiceType) Type {
	switch input {
	case entities.ChoiceTypeProficiency:
		return TypeProficiency
	case entities.ChoiceTypeLanguage:
		return TypeLanguage
	case entities.ChoiceTypeEquipment:
		return TypeEquipment
	default:
		return TypeUnset
	}
}

func typeToChoiceType(input Type) entities.ChoiceType {
	switch input {
	case TypeProficiency:
		return entities.ChoiceTypeProficiency
	case TypeLanguage:
		return entities.ChoiceTypeLanguage
	case TypeEquipment:
		return entities.ChoiceTypeEquipment
	default:
		return entities.ChoiceTypeUnset
	}
}
