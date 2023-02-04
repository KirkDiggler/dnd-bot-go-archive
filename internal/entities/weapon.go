package entities

import (
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

func (w *Weapon) GetEquipmentType() string {
	return "weapon"
}
