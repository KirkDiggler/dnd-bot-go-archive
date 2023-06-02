package attack

import (
	"fmt"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/damage"
)

type Result struct {
	AttackRoll   int
	AttackType   damage.Type
	DamageRoll   int
	AttackResult *dice.RollResult
	DamageResult *dice.RollResult
}

func (r *Result) String() string {
	return fmt.Sprintf("attack: %d, type: %s, damage: %d", r.AttackRoll, r.AttackType, r.DamageRoll)
}

func RollAttack(attackBonus int, dmg *damage.Damage) (*Result, error) {
	attackResult, err := dice.Roll(1, 20, attackBonus)
	if err != nil {
		return nil, err
	}

	dmgResult, err := dice.Roll(dmg.DiceCount, dmg.DiceSize, dmg.Bonus)
	if err != nil {
		return nil, err
	}
	dmgValue := dmgResult.Total
	attackRoll := attackResult.Total
	switch attackResult.Total {
	case 20:
		critDmg, err := dice.Roll(dmg.DiceCount, dmg.DiceSize, dmg.Bonus)
		if err != nil {
			return nil, err
		}

		dmgValue = dmgValue + critDmg.Total
	case 1:
		attackRoll = 0
	}

	return &Result{
		AttackRoll:   attackRoll,
		AttackType:   dmg.DamageType,
		DamageRoll:   dmgValue,
		AttackResult: attackResult,
		DamageResult: dmgResult,
	}, nil
}
