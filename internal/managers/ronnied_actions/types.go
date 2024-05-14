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

type AssignDrinkInput struct {
	GameID        string
	PlayerID      string
	SessionRollID string
	AssignedTo    string
}

type AssignDrinkOutput struct{}

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

type CreateSessionRollInput struct {
	SessionID    string
	Type         ronnied.RollType
	Participants []string
}

type CreateSessionRollOutput struct {
	SessionRoll *ronnied.SessionRoll
}

type GetSessionRollInput struct {
	SessionRollID string
}

type GetSessionRollOutput struct {
	SessionRoll *ronnied.SessionRoll
}

type UpdateSessionInput struct {
	Session *ronnied.Session
}

type UpdateSessionOutput struct {
	Session *ronnied.Session
}

type JoinSessionInput struct {
	SessionID     string
	SessionRollID string
	PlayerID      string
}

type JoinSessionOutput struct {
	Session     *ronnied.Session
	SessionRoll *ronnied.SessionRoll
}

type AddSessionRollInput struct {
	SessionRollID string
	PlayerID      string
}

type AddSessionRollOutput struct {
	SessionEntry *ronnied.SessionEntry
}

type UpdateSessionRollInput struct {
	SessionRoll *ronnied.SessionRoll
}

type UpdateSessionRollOutput struct {
	SessionRoll *ronnied.SessionRoll
}

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	Session *ronnied.Session
}
