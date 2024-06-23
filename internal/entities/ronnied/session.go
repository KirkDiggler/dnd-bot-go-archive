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
	ID   string
	Name string
}
type Session struct {
	ID             string
	MessageID      string // TODO: move to SessionRoll
	SessionDate    *time.Time
	GameID         string
	Players        []*Player
	StartedDate    *time.Time
	SessionRollIDs []string
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
	if len(sr.Entries) != len(sr.Players) {
		return false
	}

	for _, entry := range sr.Entries {
		if !entry.IsComplete() {
			return false
		}
	}

	return true
}

func (sr *SessionRoll) GetLosers() []*SessionEntry {
	lowestRoll := 6
	for _, entry := range sr.Entries {
		if entry.Roll < lowestRoll {
			lowestRoll = entry.Roll
		}
	}

	var losers []*SessionEntry
	for _, entry := range sr.Entries {
		if entry.Roll == lowestRoll {
			losers = append(losers, entry)
		}
	}

	return losers
}

func (sr *SessionRoll) GetWinners() []*SessionEntry {
	highestRoll := 0
	for _, entry := range sr.Entries {
		if entry.Roll > highestRoll {
			highestRoll = entry.Roll
		}
	}

	var winners []*SessionEntry
	for _, entry := range sr.Entries {
		if entry.Roll == highestRoll {
			winners = append(winners, entry)
		}
	}

	return winners
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

func (se *SessionEntry) IsComplete() bool {
	if se.Roll == 0 {
		return false
	}

	if se.Roll == 6 && se.AssignedTo == "" {
		return false
	}

	return true
}
func (se *SessionEntry) String() string {
	return fmt.Sprintf("Roll: %d", se.Roll)
}
