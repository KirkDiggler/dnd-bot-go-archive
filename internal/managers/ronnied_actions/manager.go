package ronnied_actions

import (
	"context"
	"errors"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/game"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/session"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"math/rand"
)

type Manager struct {
	gameRepo    game.Interface
	sessionRepo session.Interface
}

type ManagerConfig struct {
	GameRepo    game.Interface
	SessionRepo session.Interface
}

func NewManager(cfg *ManagerConfig) (*Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.GameRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.GameRepo")
	}

	if cfg.SessionRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.SessionRepo")
	}

	return &Manager{
		gameRepo:    cfg.GameRepo,
		sessionRepo: cfg.SessionRepo,
	}, nil
}

func (m *Manager) UpdateSession(ctx context.Context, input *UpdateSessionInput) (*UpdateSessionOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.Session == nil {
		return nil, dnderr.NewMissingParameterError("input.Session")
	}

	if input.Session.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.Session.GameID")
	}

	_, err := m.sessionRepo.Update(ctx, &session.UpdateInput{
		Session: input.Session,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateSessionOutput{
		Session: input.Session,
	}, nil
}

func (m *Manager) GetSessionRoll(ctx context.Context, input *GetSessionRollInput) (*GetSessionRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionRollID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionRollID")
	}

	result, err := m.sessionRepo.GetSessionRoll(ctx, &session.GetSessionRollInput{
		ID: input.SessionRollID,
	})
	if err != nil {
		return nil, err
	}

	return &GetSessionRollOutput{
		SessionRoll: result.SessionRoll,
	}, nil
}

func (m *Manager) AddSessionRoll(ctx context.Context, input *AddSessionRollInput) (*AddSessionRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionRollID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionRollID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	roll := rand.Intn(6) + 1

	result, err := m.sessionRepo.AddEntry(ctx, &session.AddEntryInput{
		SessionRollID: input.SessionRollID,
		PlayerID:      input.PlayerID,
		Roll:          roll,
	})
	if err != nil {
		return nil, err
	}

	slog.Info("AddSessionRoll", "result", result)

	return &AddSessionRollOutput{
		SessionEntry: result.SessionEntry,
	}, nil
}

func (m *Manager) CreateSession(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	// check that this game exists
	_, err := m.gameRepo.Get(ctx, &game.GetInput{
		ID: input.GameID,
	})
	if err != nil {
		if errors.Is(err, internal.ErrRecordNotFound) {
			_, err = m.gameRepo.Create(ctx, &game.CreateInput{
				Game: &ronnied.Game{
					ID:   input.GameID,
					Name: input.GameID,
				}})
			if err != nil {
				return nil, err
			}
		}

		return nil, err
	}

	// create a session
	result, err := m.sessionRepo.Create(ctx, &session.CreateInput{
		GameID: input.GameID,
	})
	if err != nil {
		return nil, err
	}

	return &CreateSessionOutput{
		Session: result.Session,
	}, nil
}

func (m *Manager) CreateSessionRoll(ctx context.Context, input *CreateSessionRollInput) (*CreateSessionRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionID")
	}

	if input.Participants == nil {
		return nil, dnderr.NewMissingParameterError("input.Players")
	}

	rollResult, err := m.sessionRepo.CreateRoll(ctx, &session.CreateRollInput{
		SessionID:    input.SessionID,
		Type:         input.Type,
		Participants: input.Participants,
	})
	if err != nil {
		return nil, err
	}

	return &CreateSessionRollOutput{
		SessionRoll: rollResult.SessionRoll,
	}, nil
}

func (m *Manager) JoinSession(ctx context.Context, input *JoinSessionInput) (*JoinSessionOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	// join the session
	result, err := m.sessionRepo.Join(ctx, &session.JoinInput{
		SessionID: input.SessionID,
		PlayerID:  input.PlayerID,
	})
	if err != nil {
		return nil, err
	}

	_, err = m.sessionRepo.JoinSessionRoll(ctx, &session.JoinSessionRollInput{
		SessionRollID: input.SessionRollID,
		PlayerID:      input.PlayerID,
	})
	if err != nil {
		return nil, err
	}

	return &JoinSessionOutput{
		Session: result.Session,
	}, nil
}

func (m *Manager) GetSession(ctx context.Context, input *GetSessionInput) (*GetSessionOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	result, err := m.sessionRepo.Get(ctx, &session.GetInput{
		ID: input.SessionID,
	})
	if err != nil {
		return nil, err
	}

	return &GetSessionOutput{
		Session: result.Session,
	}, nil
}

func (m *Manager) CreateGame(ctx context.Context, input *CreateGameInput) (*CreateGameOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.Game == nil {
		return nil, dnderr.NewMissingParameterError("input.Game")
	}

	if input.Game.Name == "" {
		return nil, dnderr.NewMissingParameterError("input.Game.Name")
	}

	result, err := m.gameRepo.Create(ctx, &game.CreateInput{
		Game: input.Game,
	})
	if err != nil {
		return nil, err
	}

	return &CreateGameOutput{
		Game: result.Game,
	}, nil
}

func (m *Manager) getOrCreateGame(ctx context.Context, input *JoinGameInput) (*ronnied.Game, error) {
	gameResult, err := m.gameRepo.Get(ctx, &game.GetInput{
		ID: input.GameID,
	})
	if err != nil {
		if errors.Is(err, redis.Nil) { // TODO: make const error types (e.g. ErrNotFound)
			createResult, createErr := m.CreateGame(ctx, &CreateGameInput{
				Game: &ronnied.Game{
					ID:      input.GameID,
					Name:    input.GameName,
					Players: []string{},
				},
			})
			if createErr != nil {
				return nil, createErr
			}

			return createResult.Game, nil
		}

		return nil, err
	}

	return gameResult.Game, nil
}

func (m *Manager) JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	result, err := m.getOrCreateGame(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Players == nil {
		result.Players = []string{}
	}

	for _, player := range result.Players {
		if player == input.PlayerID {
			return nil, dnderr.NewInvalidEntityError("player is already in the game")
		}
	}

	result.Players = append(result.Players, input.PlayerID)

	_, err = m.gameRepo.Create(ctx, &game.CreateInput{
		Game: result,
	})
	if err != nil {
		return nil, err
	}

	_, err = m.gameRepo.Join(ctx, &game.JoinInput{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
	})
	if err != nil {
		return nil, err
	}

	return &JoinGameOutput{}, nil
}

func (m *Manager) ListTabs(ctx context.Context, input *ListTabsInput) (*ListTabsOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	gameResult, err := m.gameRepo.Get(ctx, &game.GetInput{
		ID: input.GameID,
	})
	if err != nil {
		return nil, err
	}

	tabs := make([]*ronnied.Tab, 0)
	for _, player := range gameResult.Game.Players {
		tabCount, err := m.GetTab(ctx, &GetTabInput{
			GameID:   input.GameID,
			PlayerID: player,
		})
		if err != nil {
			return nil, err
		}

		tabs = append(tabs, &ronnied.Tab{
			PlayerID: player,
			Count:    tabCount.Count,
		})
	}

	return &ListTabsOutput{
		Tabs: tabs,
	}, nil
}

func (m *Manager) AddRolls(ctx context.Context, input *AddRollsInput) (*AddRollsOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	if input.RollCount <= 0 {
		return nil, dnderr.NewInvalidParameterError("input.RollCount must be greater than 0. Your value was:", input.RollCount)
	}

	rolls := make([]int, input.RollCount)
	for i := 0; i < input.RollCount; i++ {
		rolls[i] = rand.Intn(6) + 1
	}

	results := make([]*RollResult, len(rolls))
	for i, roll := range rolls {
		addRollOutput, err := m.AddRoll(ctx, &AddRollInput{
			GameID:   input.GameID,
			PlayerID: input.PlayerID,
			Roll:     roll,
		})
		if err != nil {
			return nil, err
		}

		rollResult := &RollResult{
			PlayerID: input.PlayerID,
			Roll:     roll,
		}

		if addRollOutput.Success {
			rollResult.AssignedTo = addRollOutput.AssignedTo
		}

		results[i] = rollResult
	}

	return &AddRollsOutput{
		Results: results,
		Success: len(results) > 0, // this cant happen rig
	}, nil
}

func (m *Manager) AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}
	gameResult, err := m.gameRepo.Get(ctx, &game.GetInput{
		ID: input.GameID,
	})
	if err != nil {
		slog.Info("gameRepo.Get", "error", err)

		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	if !gameResult.Game.HasPlayer(input.PlayerID) {
		return &AddRollOutput{}, nil
	}

	// TODO: check if the player is in the game
	// Make sure there is 1 other player in the game
	if shouldAddEntry(input.Roll) {
		var assignedTo = input.PlayerID

		if input.Roll == 6 {
			if gameResult.Game.Players == nil || len(gameResult.Game.Players) < 2 {
				return nil, dnderr.NewInvalidEntityError("game must have at least 2 players")
			}

			// select a random membership for the game
			availableMemberships := make([]string, 0)
			for _, membership := range gameResult.Game.Players {
				if membership != input.PlayerID {
					availableMemberships = append(availableMemberships, membership)
				}
			}

			randIndex := rand.Intn(len(availableMemberships))
			assignedTo = availableMemberships[randIndex]

		}

		_, err := m.gameRepo.AddEntry(ctx, &game.AddEntryInput{
			GameID:     input.GameID,
			PlayerID:   input.PlayerID,
			Roll:       input.Roll,
			AssignedTo: assignedTo,
		})
		if err != nil {
			slog.Info("gameRepo.AddEntry", "error", err)

			return nil, fmt.Errorf("failed to add entry: %w", err)
		}

		return &AddRollOutput{
			AssignedTo: assignedTo,
			Success:    true,
		}, nil
	}

	return &AddRollOutput{
		Success: false,
	}, nil
}

func (m *Manager) GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	result, err := m.gameRepo.GetTab(ctx, &game.GetTabInput{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
	})
	if err != nil {
		if errors.Is(err, internal.ErrRecordNotFound) {
			return nil, internal.ErrTabEmpty
		}

		return nil, err
	}

	return &GetTabOutput{
		Count: result.Count,
	}, nil
}

func (m *Manager) PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	_, err := m.gameRepo.PayDrink(ctx, &game.PayDrinkInput{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
	})
	if err != nil {
		if errors.Is(err, internal.ErrRecordNotFound) {
			return nil, internal.ErrTabEmpty
		}

		return nil, err
	}

	tab := &GetTabInput{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
	}

	tabResult, err := m.GetTab(ctx, tab)
	if err != nil {
		return nil, err
	}

	return &PayDrinkOutput{
		TabRemaining: tabResult.Count,
	}, nil

}

func shouldAddEntry(roll int) bool {
	return roll == 1 || roll == 6
}
