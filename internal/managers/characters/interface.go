package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	AddProficiency(ctx context.Context, char *entities.Character, reference *entities.ReferenceItem) (*entities.Character, error)
	Put(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*entities.Character, error)
	List(ctx context.Context, ownerID string) ([]*entities.Character, error)
	GetChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType) ([]*entities.Choice, error)
	SaveChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType, choices []*entities.Choice) error
	SaveState(ctx context.Context, state *entities.CharacterCreation) (*entities.CharacterCreation, error)
	GetState(ctx context.Context, id string) (*entities.CharacterCreation, error)
	AddInventory(ctx context.Context, char *entities.Character, key string) (*entities.Character, error)
	CreateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	UpdateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	GetEncounter(ctx context.Context, id string) (*entities.Encounter, error)
}
