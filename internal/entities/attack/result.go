package attack

import (
	"fmt"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/damage"
)

type Result struct {
	AttackRoll int
	AttackType damage.Type
	DamageRoll int
}

func (r *Result) String() string {
	return fmt.Sprintf("attack: %d, type: %s, damage: %d", r.AttackRoll, r.AttackType, r.DamageRoll)
}

func RollAttack(attackBonus, damageBonus int, dmg *damage.Damage) (*Result, error) {
	attackResult, err := dice.Roll(1, 20)
	if err != nil {
		return nil, err
	}

	dmgResult, err := dice.Roll(dmg.DiceCount, dmg.DiceSize)
	if err != nil {
		return nil, err
	}

	return &Result{
		AttackRoll: attackBonus + attackResult.Total,
		AttackType: dmg.DamageType,
		DamageRoll: damageBonus + dmgResult.Total,
	}, nil
}
