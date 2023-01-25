package choice

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, input *GetInput) (*GetOutput, error)
	Put(ctx context.Context, input *PutInput) error
}
