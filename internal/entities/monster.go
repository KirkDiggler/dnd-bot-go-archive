package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/attack"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/damage"
	"math/rand"
)

type Monster struct {
	ID          string           `json:"id"`
	Template    *MonsterTemplate `json:"template"`
	CharacterID string           `json:"character_id"`
	CurrentHP   int              `json:"current_hp"`
	Key         string           `json:"key"`
}

// Attack selects a random action and performs the attack
//
// returns an ampty result if the monster has no actions
func (m *Monster) Attack() ([]*attack.Result, error) {
	if len(m.Template.Actions) == 0 {
		return []*attack.Result{}, nil
	}

	randRoll := rand.Intn(len(m.Template.Actions) - 1)
	action := m.Template.Actions[randRoll]

	results := make([]*attack.Result, 0, len(action.Damage))
	for _, dmg := range action.Damage {
		attackResult, err := attack.RollAttack(action.AttackBonus, 0, dmg)
		if err != nil {
			return nil, err
		}
		results = append(results, attackResult)
	}
	
	return results, nil
}

type MonsterTemplate struct {
	Key             string           `json:"key"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	ArmorClass      int              `json:"armor_class"`
	HitPoints       int              `json:"hit_points"`
	HitDice         string           `json:"hit_dice"`
	Actions         []*MonsterAction `json:"actions"`
	XP              int              `json:"xp"`
	ChallengeRating float32          `json:"challenge_rating"`
}

type MonsterAction struct {
	Name        string           `json:"name"`
	AttackBonus int              `json:"attack_bonus"`
	Description string           `json:"desc"`
	Damage      []*damage.Damage `json:"damage"`
}
