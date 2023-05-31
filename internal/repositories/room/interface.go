package room

import (
	"context"
)

type Repository interface {
	CreateRoom(ctx context.Context, room *Data) (*Data, error)
	UpdateRoom(ctx context.Context, room *Data) (*Data, error)
	GetRoom(ctx context.Context, id string) (*Data, error)
	ListRooms(ctx context.Context, owner string) ([]*Data, error)
}
