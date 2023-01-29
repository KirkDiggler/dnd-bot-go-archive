package entities

import "fmt"

type Attribute string

const (
	AttributeStrength     Attribute = "str"
	AttributeDexterity    Attribute = "dex"
	AttributeConstitution Attribute = "con"
	AttributeIntelligence Attribute = "int"
	AttributeWisdom       Attribute = "wis"
	AttributeCharisma     Attribute = "cha"
)

type AbilityScore struct {
	Score int
	Bonus int
}

func (a *AbilityScore) AddBonus(bonus int) *AbilityScore {
	a.Bonus += bonus

	return a
}

func (a *AbilityScore) Display() string {
	return fmt.Sprintf("%d (%+d)", a.Score, a.Bonus)
}
