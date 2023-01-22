package character

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Put(ctx context.Context, character *Data) (*Data, error) {
	args := m.Called(ctx, character)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), nil
}

func (m *Mock) Get(ctx context.Context, id string) (*Data, error) {
	args := m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), nil
}
