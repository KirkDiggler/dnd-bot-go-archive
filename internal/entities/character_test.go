package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCharacter_ToData(t *testing.T) {
	char := &Character{
		ID:      "id",
		Name:    "name",
		OwnerID: "ownerID",
		Race: &Race{
			Key: "race-key",
		},
		Class: &Class{
			Key: "class-key",
		},
		Attribues: map[Attribute]*AbilityScore{
			AttributeStrength:     {Score: 16},
			AttributeDexterity:    {Score: 15},
			AttributeConstitution: {Score: 14},
			AttributeIntelligence: {Score: 13},
			AttributeWisdom:       {Score: 12},
			AttributeCharisma:     {Score: 11},
		},
	}

	expected := &character.Data{
		ID:       "id",
		Name:     "name",
		OwnerID:  "ownerID",
		RaceKey:  "race-key",
		ClassKey: "class-key",
		Attributes: &character.AttributeData{
			Str: 16,
			Dex: 15,
			Con: 14,
			Int: 13,
			Wis: 12,
			Cha: 11,
		},
	}

	actual := char.ToData()
	assert.Equal(t, expected, actual)
}
