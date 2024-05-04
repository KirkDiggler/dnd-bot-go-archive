package ronnied

import "time"

type RollType string

const (
	RollTypeStart   RollType = "start"
	RollTypeRollOff RollType = "roll_off"
)

type Session struct {
	ID          string
	SessionDate *time.Time
	GameID      string
	Players     []string
	StartedDate *time.Time
}

func (s *Session) HasPlayer(playerID string) bool {
	for _, player := range s.Players {
		if player == playerID {
			return true
		}
	}

	return false
}

type SessionRoll struct {
	ID           string
	SessionID    string
	Type         RollType
	Participants []string // players involved in this roll
	Entries      []*SessionEntry
}

type SessionEntry struct {
	ID         string
	SessionID  string
	PlayerID   string
	Roll       int
	AssignedTo string
}
