package entities

import "fmt"

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
	OwnerID     string
	LastToken   string
	Step        CreateStep
	Steps       CreateStep
}

func (c *CharacterCreation) String() string {
	return fmt.Sprintf("CharacterCreation{CharacterID: %s, OwnerID: %s, LastToken: %s, Step: %d, Steps: %d}", c.CharacterID, c.OwnerID, c.LastToken, c.Step, c.Steps)
}
