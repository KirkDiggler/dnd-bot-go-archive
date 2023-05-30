package monster

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetMonster(ctx context.Context, key string) (*entities.Monster, error) {
	args := m.Called(ctx, key)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Monster), nil
}

func (m *Mock) PutMonster(ctx context.Context, monster *entities.Monster) error {
	args := m.Called(ctx, monster)

	return args.Error(0)
}
