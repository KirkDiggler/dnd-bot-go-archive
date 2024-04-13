package ronnied_actions

import (
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
)

type CreateGameInput struct {
	Game *ronnied.Game
}

type CreateGameOutput struct {
	Game *ronnied.Game
}

type JoinGameInput struct {
	GameID   string
	GameName string
	PlayerID string
}

type JoinGameOutput struct {
}

type AddRollInput struct {
	GameID   string
	PlayerID string
	Roll     int
}

type AddRollOutput struct {
	AssignedTo string
	Success    bool
}

type GetTabInput struct {
	GameID   string
	PlayerID string
}

type GetTabOutput struct {
	Count int
}

type PayDrinkInput struct {
	GameID   string
	PlayerID string
}

type PayDrinkOutput struct{}

type ListTabsInput struct {
	GameID string
}

type ListTabsOutput struct {
	Tabs []*ronnied.Tab
}