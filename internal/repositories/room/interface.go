package room

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, room *Data) (*Data, error)
	Update(ctx context.Context, room *Data) (*Data, error)
	Get(ctx context.Context, id string) (*Data, error)
	ListByPlayer(ctx context.Context, input *ListByPlayerInput) ([]*Data, error)
}
