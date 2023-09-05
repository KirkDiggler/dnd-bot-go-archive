package encounter

import "context"

type Repository interface {
	Create(ctx context.Context, encounter *Data) (*Data, error)
	Update(ctx context.Context, encounter *Data) (*Data, error)
	Get(ctx context.Context, id string) (*Data, error)
	ListByPlayer(ctx context.Context, playerID string) ([]*Data, error)
}
