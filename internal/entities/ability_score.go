package entities

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
