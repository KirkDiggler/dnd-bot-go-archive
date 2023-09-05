package entities

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/damage"
	"github.com/stretchr/testify/suite"
	"testing"
)

type suiteMonster struct {
	suite.Suite

	ctx     context.Context
	monster *Monster
}

func (s *suiteMonster) SetupTest() {
	s.ctx = context.Background()
	s.monster = &Monster{
		ID: "Goblin-ID",
		Template: &MonsterTemplate{
			Name: "Goblin",
			Actions: []*MonsterAction{
				{
					Name:        "Scimitar",
					AttackBonus: 4,
					Description: "Melee Weapon Attack: +4 to hit, reach 5 ft., one target. Hit: 5 (1d6 + 2) slashing damage.",
					Damage: []*damage.Damage{
						{
							DiceCount:  1,
							DiceSize:   6,
							Bonus:      2,
							DamageType: damage.TypeSlashing,
						},
					},
				},
			},
		},
	}
}

func (s *suiteMonster) TestAttack() {
	results, err := s.monster.Attack()
	s.NoError(err)
	s.Len(results, 1)
	s.True(results[0].AttackRoll > 0)
	s.True(results[0].DamageRoll > 0)
	s.Equal(results[0].DamageRoll, results[0].DamageResult.Total)
	s.Equal(results[0].AttackRoll, results[0].AttackResult.Total)
	s.Equal(results[0].DamageRoll, results[0].DamageResult.Rolls[0]+results[0].DamageResult.Bonus)
	s.Equal(results[0].AttackRoll, results[0].AttackResult.Rolls[0]+results[0].AttackResult.Bonus)
}

func TestMonster(t *testing.T) {
	suite.Run(t, new(suiteMonster))
}
