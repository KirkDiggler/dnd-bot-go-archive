package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	Put(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*entities.Character, error)
	GetChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType) ([]*entities.Choice, error)
	SaveChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType, choices []*entities.Choice) error
	SaveState(ctx context.Context, state *entities.CharacterCreation) (*entities.CharacterCreation, error)
	GetState(ctx context.Context, id string) (*entities.CharacterCreation, error)
}
