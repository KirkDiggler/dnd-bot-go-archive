package rooms

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	LoadRoom(ctx context.Context, input *LoadRoomInput) (*LoadRoomOutput, error)
	HasActiveRoom(ctx context.Context, input *HasActiveRoomInput) (*HasActiveRoomOutput, error)
	Attack(ctx context.Context, playerID string) (string, error)
}

type LoadRoomInput struct {
	PlayerID string
}

type LoadRoomOutput struct {
	Room *entities.Room
}

type HasActiveRoomInput struct {
	PlayerID string
}

type HasActiveRoomOutput struct {
	HasActiveRoom bool
}
