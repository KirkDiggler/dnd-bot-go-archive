package entities

type Equipment interface {
	GetEquipmentType() string
	GetName() string
	GetKey() string
	GetSlot() Slot
}
