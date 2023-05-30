package room

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Repository interface {
	CreateRoom(ctx context.Context, room *entities.Room) (*entities.Room, error)
	UpdateRoom(ctx context.Context, room *entities.Room) (*entities.Room, error)
	GetRoom(ctx context.Context, id string) (*entities.Room, error)
	ListRooms(ctx context.Context, owner string) ([]*entities.Room, error)
}
