package character

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Put(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*entities.Character, error)
	ListByOwner(ctx context.Context, ownerID string) ([]*entities.Character, error)
	Delete(ctx context.Context, id string) error
	// Add a new method to list characters by owner and status
	ListByOwnerAndStatus(ctx context.Context, ownerID string, status ...entities.CharacterStatus) ([]*entities.Character, error)
}
