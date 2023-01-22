package character

import (
	"context"
)

type Repository interface {
	Put(ctx context.Context, character *Data) (*Data, error)
	Get(ctx context.Context, id string) (*Data, error)
}
