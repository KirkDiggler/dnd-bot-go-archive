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

type RollResult struct {
	PlayerID   string
	AssignedTo string
	Roll       int
}

type AddRollsInput struct {
	GameID    string
	PlayerID  string
	RollCount int
}

type AddRollsOutput struct {
	Results []*RollResult
	Success bool
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

type PayDrinkOutput struct {
	TabRemaining int
}

type ListTabsInput struct {
	GameID string
}

type ListTabsOutput struct {
	Tabs []*ronnied.Tab
}

type CreateSessionInput struct {
	GameID string
}

type CreateSessionOutput struct {
	Session *ronnied.Session
}

type JoinSessionInput struct {
	SessionID string
	PlayerID  string
}

type JoinSessionOutput struct {
	Session *ronnied.Session
}

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	Session *ronnied.Session
}
