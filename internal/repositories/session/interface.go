package session

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	Create(ctx context.Context, session *entities.Session) (*entities.Session, error)
	Get(ctx context.Context, userID string) (*entities.Session, error)
	Update(ctx context.Context, session *entities.Session) (*entities.Session, error)
	Delete(ctx context.Context, userID string) error
}
