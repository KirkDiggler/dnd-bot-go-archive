package entities

type Race struct {
	Key                        string           `json:"key"`
	Name                       string           `json:"name"`
	StartingProficiencyOptions *Choice          `json:"proficiency_choices"`
	StartingProficiencies      []*ReferenceItem `json:"proficiencies"`
	AbilityBonuses             []*AbilityBonus  `json:"ability_bonuses"`
}
