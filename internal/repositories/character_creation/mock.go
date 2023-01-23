package character_creation

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (_m *Mock) Get(ctx context.Context, id string) (*entities.CharacterCreation, error) {
	args := _m.Called(ctx, id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.CharacterCreation), nil
}

func (_m *Mock) Put(ctx context.Context, character *entities.CharacterCreation) (*entities.CharacterCreation, error) {
	args := _m.Called(ctx, character)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.CharacterCreation), nil
}
