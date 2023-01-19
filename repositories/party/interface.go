package party

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/entities"
)

type Interface interface {
	CreateParty(ctx context.Context, party *entities.Party) (*entities.Party, error)
	GetParty(ctx context.Context, token string) (*entities.Party, error)
}
