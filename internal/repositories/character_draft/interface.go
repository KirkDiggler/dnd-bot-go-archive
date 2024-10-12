package character_draft

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, draft *entities.CharacterDraft) (*entities.CharacterDraft, error)
	Get(ctx context.Context, id string) (*entities.CharacterDraft, error)
	Update(ctx context.Context, draft *entities.CharacterDraft) (*entities.CharacterDraft, error)
	ListByOwner(ctx context.Context, ownerID string) ([]*entities.CharacterDraft, error)
	Delete(ctx context.Context, id string) error
}
