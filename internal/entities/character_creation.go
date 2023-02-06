package entities

type CreateStep int

const (
	CreateStepSelect CreateStep = iota
	CreateStepRoll
	CreateStepProficiency
	CreateStepEquipment
	CreateStepEquipCharacter
	CreateStepDone
)

type CharacterCreation struct {
	CharacterID string
	LastToken   string
	Step        CreateStep
}
