package entities

type ArmorClass struct {
	Base     int  `json:"armor_class"`
	DexBonus bool `json:"dex_bonus"`
	MaxBonus int  `json:"max_bonus"`
}

type Armor struct {
	Base                BasicEquipment `json:"base"`
	ArmorClass          *ArmorClass    `json:"armor_class"`
	StrMin              int            `json:"str_minimum"`
	StealthDisadvantage bool           `json:"stealth_disadvantage"`
}

func (e *Armor) GetEquipmentType() string {
	return "Armor"
}

func (e *Armor) GetName() string {
	return e.Base.Name
}

func (e *Armor) GetKey() string {
	return e.Base.Key
}
