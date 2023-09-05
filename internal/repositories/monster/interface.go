package monster

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Interface interface {
	GetMonster(ctx context.Context, key string) (*entities.Monster, error)
	PutMonster(ctx context.Context, monster *entities.Monster) (*entities.Monster, error)
}
