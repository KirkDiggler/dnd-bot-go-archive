package entities

type BasicEquipment struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Cost   *Cost  `json:"cost"`
	Weight int    `json:"weight"`
}

func (e *BasicEquipment) GetEquipmentType() string {
	return "basic_equipment"
}
