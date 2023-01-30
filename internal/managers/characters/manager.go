package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/choice"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_creation"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

type manager struct {
	charRepo   character.Repository
	stateRepo  character_creation.Repository
	choiceRepo choice.Repository
	client     dnd5e.Client
}

type Config struct {
	CharacterRepo character.Repository
	StateRepo     character_creation.Repository
	ChoiceRepo    choice.Repository
	Client        dnd5e.Client
}

func New(cfg *Config) (Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.CharacterRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharacterRepo")
	}

	if cfg.StateRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.StateRepo")
	}

	if cfg.ChoiceRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.ChoiceRepo")
	}

	return &manager{
		charRepo:   cfg.CharacterRepo,
		stateRepo:  cfg.StateRepo,
		choiceRepo: cfg.ChoiceRepo,
		client:     cfg.Client,
	}, nil
}

func (m *manager) AddProficiency(ctx context.Context, char *entities.Character, reference *entities.ReferenceItem) (*entities.Character, error) {
	if char == nil {
		return nil, dnderr.NewMissingParameterError("char")
	}

	if reference == nil {
		return nil, dnderr.NewMissingParameterError("reference")
	}

	if reference.Type != entities.ReferenceTypeProficiency {
		return nil, dnderr.NewInvalidParameterError("reference.Type", "must be a proficiency")
	}

	proficiency, err := m.client.GetProficiency(reference.Key)
	if err != nil {
		return nil, err
	}

	char.AddProficiency(proficiency)

	return char, nil
}

func (m *manager) SaveChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType, choices []*entities.Choice) error {
	if characterID == "" {
		return dnderr.NewMissingParameterError("characterID")
	}

	if choiceType == "" {
		return dnderr.NewMissingParameterError("choiceType")
	}

	if choices == nil {
		return dnderr.NewMissingParameterError("choices")
	}

	if len(choices) == 0 {
		return dnderr.NewMissingParameterError("choices")
	}

	return m.choiceRepo.Put(ctx, &choice.PutInput{
		CharacterID: characterID,
		Type:        choiceType,
		Choices:     choices,
	})
}

func (m *manager) GetChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType) ([]*entities.Choice, error) {
	if characterID == "" {
		return nil, dnderr.NewMissingParameterError("characterID")
	}

	if choiceType == "" {
		return nil, dnderr.NewMissingParameterError("choiceType")
	}

	data, err := m.choiceRepo.Get(ctx, &choice.GetInput{
		CharacterID: characterID,
		Type:        choiceType,
	})
	if err != nil {
		return nil, err
	}

	return data.Choices, nil
}

func (m *manager) Put(ctx context.Context, character *entities.Character) (*entities.Character, error) {
	if character == nil {
		return nil, dnderr.NewMissingParameterError("character")
	}

	if character.Name == "" {
		return nil, dnderr.NewMissingParameterError("character.Name")
	}

	if character.OwnerID == "" {
		return nil, dnderr.NewMissingParameterError("character.OwnerID")
	}

	if character.Race == nil {
		return nil, dnderr.NewMissingParameterError("character.Race")
	}

	if character.Class == nil {
		return nil, dnderr.NewMissingParameterError("character.Class")
	}

	data, err := m.charRepo.Put(ctx, character)
	if err != nil {
		return nil, err
	}

	character.ID = data.ID

	return character, nil
}

func (m *manager) Get(ctx context.Context, id string) (*entities.Character, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	data, err := m.charRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return m.characterFromData(ctx, data)
}

// TODO: Move to state manager
func (m *manager) SaveState(ctx context.Context, state *entities.CharacterCreation) (*entities.CharacterCreation, error) {
	if state == nil {
		return nil, dnderr.NewMissingParameterError("state")
	}

	if state.CharacterID == "" {
		return nil, dnderr.NewMissingParameterError("state.CharacterID")
	}

	result, err := m.stateRepo.Put(ctx, state)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) GetState(ctx context.Context, characterID string) (*entities.CharacterCreation, error) {
	if characterID == "" {
		return nil, dnderr.NewMissingParameterError("characterID")
	}

	result, err := m.stateRepo.Get(ctx, characterID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
