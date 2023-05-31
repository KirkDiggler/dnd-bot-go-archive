package room

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateRoom(ctx context.Context, room *Data) (*Data, error) {
	args := m.Called(ctx, room)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) GetRoom(ctx context.Context, id string) (*Data, error) {
	args := m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) UpdateRoom(ctx context.Context, room *Data) (*Data, error) {
	args := m.Called(ctx, room)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Data), args.Error(1)
}

func (m *Mock) ListRooms(ctx context.Context, owner string) ([]*Data, error) {
	args := m.Called(ctx, owner)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Data), args.Error(1)
}
