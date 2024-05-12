package ronnied

import "time"

type RollType string

const (
	RollTypeStart   RollType = "start"
	RollTypeRollOff RollType = "roll_off"
)

type Session struct {
	ID          string
	MessageID   string
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
	ID        string
	SessionID string
	Type      RollType
	Players   []string // players involved in this roll
	Entries   []*SessionEntry
}

func (sr *SessionRoll) HasPlayer(playerID string) bool {
	for _, player := range sr.Players {
		if player == playerID {
			return true
		}
	}

	return false
}

type SessionEntry struct {
	ID            string
	SessionRollID string
	PlayerID      string
	Roll          int
	AssignedTo    string
}
