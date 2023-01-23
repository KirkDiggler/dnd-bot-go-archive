package entities

type Class struct {
	Key                string                `json:"key"`
	Name               string                `json:"name"`
	ProficiencyChoices []*ProficiencyChoices `json:"proficiency_choices"`
}
