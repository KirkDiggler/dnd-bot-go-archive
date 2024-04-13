package ronnied

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type Game struct {
	ID        string     `json:"id"`   // Channel ID
	Name      string     `name:"name"` // channel name when created
	StartedAt *time.Time `json:"started_at"`
	Players   []string   `json:"players"`
}

func (g *Game) HasPlayer(playerID string) bool {
	for _, player := range g.Players {
		if player == playerID {
			return true
		}
	}

	return false
}

func (g *Game) String() string {
	return fmt.Sprintf("Game{ID: %s, Name: %s, StartedAt: %s}", g.ID, g.Name, g.StartedAt)
}

func (g *Game) MarshalGameString() string {
	output, err := json.Marshal(g)
	if err != nil {
		slog.Warn("error marshalling game", g.String(), err.Error())
		return ""
	}

	return string(output)
}

func UnmarshalGameString(input string) (*Game, error) {
	out := &Game{}
	err := json.Unmarshal([]byte(input), out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type GameEntry struct {
	ID             string     `json:"id"`
	GameID         string     `json:"game_id"`
	PlayerID       string     `json:"player_id"`
	Roll           int        `json:"roll"`
	AssignedTo     string     `json:"assigned_to"`
	Status         string     `json:"status"` // TODO: what statuses can we have?
	CreatedDate    *time.Time `json:"created_date"`
	ReconciledDate *time.Time `json:"reconciled_date"`
}

type GameTab struct {
	GameID      string `json:"game_id"`
	PlayerID    string `json:"player_id"`
	GameEntryID string `json:"game_entry_id"`
}

type Tab struct {
	PlayerID string `json:"player_id"`
	Count    int    `json:"count"`
}
