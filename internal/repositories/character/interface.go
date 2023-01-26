package character

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Put(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*Data, error)
}
