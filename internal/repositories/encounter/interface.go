package encounter

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	Update(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	Get(ctx context.Context, id string) (*entities.Encounter, error)
}
