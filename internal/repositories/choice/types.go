package choice

import "github.com/KirkDiggler/dnd-bot-go/internal/entities"

type GetInput struct {
	CharacterID string
	Type        Type
}

type GetOutput struct {
	CharacterID string
	Type        Type
	Choices     []*entities.Choice
}

type PutInput struct {
	CharacterID string
	Type        Type
	Choices     []*entities.Choice
}
