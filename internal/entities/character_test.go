package entities

import "testing"

func TestCharacter_AddAbilityBonus(t *testing.T) {
	type fields struct {
		AbilityScores map[Attribute]*AbilityScore
	}
	type args struct {
		abilityBonus  *AbilityBonus
		expectedBonus int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestCharacter_AddAbilityBonus",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeStrength:     {Score: 10},
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
					AttributeCharisma:     {Score: 10},
				},
			},
			args: args{
				abilityBonus: &AbilityBonus{
					Attribute: AttributeStrength,
					Bonus:     2,
				},
				expectedBonus: 2,
			},
		},
		{
			name: "TestCharacter_AddToExistingAbilityBonus",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeStrength:     {Score: 12, Bonus: 1},
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
					AttributeCharisma:     {Score: 10},
				},
			},
			args: args{
				abilityBonus: &AbilityBonus{
					Attribute: AttributeStrength,
					Bonus:     2,
				},
				expectedBonus: 3,
			},
		},
		{
			name: "TestCharacter_AddToNewAttributeAbilityBonus",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
					AttributeCharisma:     {Score: 10},
				},
			},
			args: args{
				abilityBonus: &AbilityBonus{
					Attribute: AttributeStrength,
					Bonus:     2,
				},
				expectedBonus: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Character{
				Attribues: tt.fields.AbilityScores,
			}

			c.AddAbilityBonus(tt.args.abilityBonus)

			if c.Attribues[tt.args.abilityBonus.Attribute].Bonus != tt.args.expectedBonus {
				t.Errorf("expected bonus to be %d, got %d", tt.args.expectedBonus, c.Attribues[tt.args.abilityBonus.Attribute].Bonus)
			}
		})
	}
}

func TestCharacter_AddAttribute(t *testing.T) {
	type fields struct {
		AbilityScores map[Attribute]*AbilityScore
	}
	type args struct {
		attribute    Attribute
		abilityScore *AbilityScore
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestCharacter_AddAttribute",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeStrength:     {Score: 0, Bonus: 2},
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
					AttributeCharisma:     {Score: 10},
				},
			},
			args: args{
				attribute: AttributeStrength,
				abilityScore: &AbilityScore{
					Score: 10,
					Bonus: 2,
				},
			},
		},
		{
			name: "TestCharacter_AddAttribute",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeStrength:     {Score: 0, Bonus: 0},
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
					AttributeCharisma:     {Score: 10},
				},
			},
			args: args{
				attribute: AttributeStrength,
				abilityScore: &AbilityScore{
					Score: 12,
					Bonus: 1,
				},
			},
		},
		{
			name: "TestCharacter_AddAttribute",
			fields: fields{
				AbilityScores: map[Attribute]*AbilityScore{
					AttributeStrength:     {Score: 10},
					AttributeDexterity:    {Score: 10},
					AttributeConstitution: {Score: 10},
					AttributeIntelligence: {Score: 10},
					AttributeWisdom:       {Score: 10},
				},
			},
			args: args{
				attribute: AttributeCharisma,
				abilityScore: &AbilityScore{
					Score: 10,
					Bonus: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Character{
				Attribues: tt.fields.AbilityScores,
			}

			c.AddAttribute(tt.args.attribute, tt.args.abilityScore.Score)
			_ = c.String()
			if _, ok := c.Attribues[tt.args.attribute]; !ok {
				t.Errorf("expected attribute %s to be present", tt.args.attribute)
			}

			if c.Attribues[tt.args.attribute].Score != tt.args.abilityScore.Score {
				t.Errorf("expected score to be %d, got %d", tt.args.abilityScore.Score, c.Attribues[tt.args.attribute].Score)
			}

			if c.Attribues[tt.args.attribute].Bonus != tt.args.abilityScore.Bonus {
				t.Errorf("expected bonus to be %d, got %d", tt.args.abilityScore.Bonus, c.Attribues[tt.args.attribute].Bonus)
			}
		})
	}
}
