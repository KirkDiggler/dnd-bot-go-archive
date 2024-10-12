package session

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	Create(ctx context.Context, userID string, draftID string) (*entities.Session, error)
	GetWithDraft(ctx context.Context, userID string) (*entities.Session, error)
	Update(ctx context.Context, session *entities.Session) (*entities.Session, error)
	UpdateLastToken(ctx context.Context, userID string, lastToken string) error
	Delete(ctx context.Context, userID string) error
}
