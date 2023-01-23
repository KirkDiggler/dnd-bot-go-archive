package entities

type Class struct {
	Key                string    `json:"key"`
	Name               string    `json:"name"`
	ProficiencyChoices []*Choice `json:"proficiency_choices"`
}
