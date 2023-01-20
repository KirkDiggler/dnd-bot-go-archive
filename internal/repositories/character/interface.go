package character

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, character *Data) (*Data, error)
	Get(ctx context.Context, id string) (*Data, error)
}
