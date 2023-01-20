package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	Create(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*entities.Character, error)
}
