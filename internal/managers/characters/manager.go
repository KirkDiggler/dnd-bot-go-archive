package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_creation"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

type manager struct {
	charRepo  character.Repository
	stateRepo character_creation.Repository
	client    dnd5e.Client
}

type Config struct {
	CharacterRepo character.Repository
	StateRepo     character_creation.Repository
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

	return &manager{
		charRepo:  cfg.CharacterRepo,
		stateRepo: cfg.StateRepo,
		client:    cfg.Client,
	}, nil
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

	data, err := m.charRepo.Put(ctx, character.ToData())
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
