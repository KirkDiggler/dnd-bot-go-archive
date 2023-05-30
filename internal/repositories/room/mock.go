package room

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateRoom(ctx context.Context, room *entities.Room) (*entities.Room, error) {
	args := m.Called(ctx, room)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Room), args.Error(1)
}

func (m *Mock) GetRoom(ctx context.Context, id string) (*entities.Room, error) {
	args := m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Room), args.Error(1)
}
