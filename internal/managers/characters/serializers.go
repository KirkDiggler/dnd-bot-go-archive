package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
)

func (m *manager) characterFromData(ctx context.Context, data *character.Data) (*entities.Character, error) {
	if data == nil {
		return nil, dnderr.NewMissingParameterError("data")
	}

	race, err := m.client.GetRace(data.RaceKey)
	if err != nil {
		return nil, err
	}

	class, err := m.client.GetClass(data.ClassKey)
	if err != nil {
		return nil, err
	}

	return &entities.Character{
		ID:      data.ID,
		Name:    data.Name,
		OwnerID: data.OwnerID,
		Race:    race,
		Class:   class,
	}, nil
}
