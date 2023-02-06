package entities

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/attack"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/damage"
)

type Weapon struct {
	Base            BasicEquipment   `json:"base"`
	Damage          *damage.Damage   `json:"damage"`
	Range           int              `json:"range"`
	WeaponCategory  string           `json:"weapon_category"`
	WeaponRange     string           `json:"weapon_range"`
	CategoryRange   string           `json:"category_range"`
	Properties      []*ReferenceItem `json:"properties"`
	TwoHandedDamage *damage.Damage   `json:"two_handed_damage"`
}

func (e *Weapon) Attack(char *Character) (*attack.Result, error) {
	var bonus int
	if e.WeaponRange == "Ranged" {
		bonus = char.Attribues[AttributeDexterity].Bonus
	} else if e.WeaponRange == "Melee" {
		bonus = char.Attribues[AttributeStrength].Bonus
	}

	// TODO: check proficiency
	if e.IsTwoHanded() {
		if e.TwoHandedDamage == nil {
			return attack.RollAttack(bonus, bonus, e.Damage)
		}

		return attack.RollAttack(bonus, bonus, e.TwoHandedDamage)
	}

	return attack.RollAttack(bonus, bonus, e.Damage)
}

func (e *Weapon) IsTwoHanded() bool {
	for _, p := range e.Properties {
		if p.Key == "two-handed" {
			return true
		}
	}

	return false
}

func (w *Weapon) GetEquipmentType() EquipmentType {
	return "Weapon"
}

func (w *Weapon) GetName() string {
	return w.Base.Name
}

func (w *Weapon) GetKey() string {
	return w.Base.Key
}

func (w *Weapon) GetSlot() Slot {
	for _, p := range w.Properties {
		if p.Key == "two-handed" {
			return SlotTwoHanded
		}
	}

	return SlotMainHand
}
