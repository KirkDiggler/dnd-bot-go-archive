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

const (
	SelectRaceStep          CreateStep = 1 << 0
	SelectClassStep         CreateStep = 1 << 1
	EnterNameStep           CreateStep = 1 << 2
	SelectBackgroundStep    CreateStep = 1 << 3
	SelectAlignmentStep     CreateStep = 1 << 4
	SelectAbilityScoresStep CreateStep = 1 << 5
	SelectSkillsStep        CreateStep = 1 << 6
	SelectEquipmentStep     CreateStep = 1 << 7
	SelectProficienciesStep CreateStep = 1 << 8
)

type CharacterCreation struct {
	CharacterID string
	LastToken   string
	Step        CreateStep
	Steps       CreateStep
}
