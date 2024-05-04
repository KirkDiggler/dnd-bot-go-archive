package session

import "github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"

type CreateInput struct {
	GameID string
}

type CreateOutput struct {
	Session *ronnied.Session
}

type UpdateInput struct {
	Session *ronnied.Session
}

type UpdateOutput struct {
	Session *ronnied.Session
}

type GetInput struct {
	ID string
}

type GetOutput struct {
	Session *ronnied.Session
}

type JoinInput struct {
	SessionID string
	PlayerID  string
}

type JoinOutput struct {
	Session *ronnied.Session
}

type CreateRollInput struct {
	SessionID    string
	Type         ronnied.RollType
	Participants []string
}

type CreateRollOutput struct {
	SessionRoll *ronnied.SessionRoll
}

type AddEntryInput struct {
	SessionID  string
	PlayerID   string
	Roll       int
	AssignedTo string
}

type AddEntryOutput struct {
	SessionEntry *ronnied.SessionEntry
}
