package character_creation

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Get(ctx context.Context, id string) (*entities.CharacterCreation, error)
	Put(ctx context.Context, character *entities.CharacterCreation) (*entities.CharacterCreation, error)
}
