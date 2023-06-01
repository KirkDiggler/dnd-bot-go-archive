package room

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Create(ctx context.Context, room *Data) (*Data, error) {
	args := m.Called(ctx, room)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) Get(ctx context.Context, id string) (*Data, error) {
	args := m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) Update(ctx context.Context, room *Data) (*Data, error) {
	args := m.Called(ctx, room)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) ListByPlayer(ctx context.Context, input *ListByPlayerInput) ([]*Data, error) {
	args := m.Called(ctx, input)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Data), args.Error(1)
}
