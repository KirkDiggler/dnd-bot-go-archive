package entities

type Attribute string

const (
	AttributeStrength     Attribute = "strength"
	AttributeDexterity    Attribute = "dexterity"
	AttributeConstitution Attribute = "constitution"
	AttributeIntelligence Attribute = "intelligence"
	AttributeWisdom       Attribute = "wisdom"
	AttributeCharisma     Attribute = "charisma"
)

type AbilityScore struct {
	Score int
	Bonus int
}
