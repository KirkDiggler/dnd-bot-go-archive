package character

import (
	"context"
)

type Repository interface {
	CreateCharacter(ctx context.Context, character *Data) (*Data, error)
	GetCharacter(ctx context.Context, id string) (*Data, error)
}
