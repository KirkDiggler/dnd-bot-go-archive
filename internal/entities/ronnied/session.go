package ronnied

import (
	"fmt"
	"time"
)

type RollType string

const (
	RollTypeStart   RollType = "start"
	RollTypeRollOff RollType = "roll_off"
)

type Player struct {
	ID    string
	MsgID string
}
type Session struct {
	ID          string
	MessageID   string // TODO: move to SessionRoll
	SessionDate *time.Time
	GameID      string
	Players     []*Player
	StartedDate *time.Time
}

func (s *Session) HasPlayer(playerID string) *Player {
	for _, player := range s.Players {
		if player.ID == playerID {
			return player
		}
	}

	return nil
}

type SessionRoll struct {
	ID        string
	SessionID string
	Type      RollType
	Players   []*Player // players involved in this roll
	Entries   []*SessionEntry
}

func (sr *SessionRoll) HasPlayer(playerID string) *Player {
	for _, player := range sr.Players {
		if player.ID == playerID {
			return player
		}
	}

	return nil
}

func (sr *SessionRoll) HasPlayerEntry(playerID string) *SessionEntry {
	for _, entry := range sr.Entries {
		if entry.PlayerID == playerID {
			return entry
		}
	}

	return nil
}

func (sr *SessionRoll) UpdatePlayerMsgID(playerID, msgID string) {
	for _, player := range sr.Players {
		if player.ID == playerID {
			player.MsgID = msgID
		}
	}
}

// LoserEntires the entries that have the lowest score
func (sr *SessionRoll) IsLoser(input *SessionEntry) bool {
	lowestRoll := 6
	for _, entry := range sr.Entries {
		if entry.Roll < lowestRoll {
			lowestRoll = entry.Roll
		}
	}

	for _, entry := range sr.Entries {
		if input.PlayerID == entry.PlayerID && entry.Roll == lowestRoll {
			return true
		}
	}

	return false
}

func (sr *SessionRoll) IsComplete() bool {
	if len(sr.Entries) == 0 {
		return false
	}

	for _, entry := range sr.Entries {
		if !entry.Completed {
			return false
		}
	}

	return true
}

type SessionEntry struct {
	ID            string
	SessionRollID string
	PlayerID      string
	Roll          int
	AssignedTo    string
	Completed     bool
}

func (se *SessionEntry) Complete() {
	se.Completed = true
}

func (se *SessionEntry) String() string {
	return fmt.Sprintf("Roll: %d", se.Roll)
}
