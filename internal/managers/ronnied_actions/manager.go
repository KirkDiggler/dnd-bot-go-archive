package ronnied_actions

import (
	"context"
	"errors"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/game"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"math/rand"
)

type Manager struct {
	gameRepo game.Interface
}

type ManagerConfig struct {
	GameRepo game.Interface
}

func NewManager(cfg *ManagerConfig) (*Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.GameRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.GameRepo")
	}

	return &Manager{
		gameRepo: cfg.GameRepo,
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

	// TODO: check if the player is in the game
	// Make sure there is 1 other player in the game
	if shouldAddEntry(input.Roll) {
		var assignedTo = input.PlayerID

		if input.Roll == 6 {
			gameResult, err := m.gameRepo.Get(ctx, &game.GetInput{
				ID: input.GameID,
			})
			if err != nil {
				slog.Info("gameRepo.Get", "error", err)

				return nil, fmt.Errorf("failed to get game: %w", err)
			}

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

	slog.Info("GetTab", "input", input)

	result, err := m.gameRepo.GetTab(ctx, &game.GetTabInput{
		GameID:   input.GameID,
		PlayerID: input.PlayerID,
	})
	if err != nil {
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
		return nil, err
	}

	return &PayDrinkOutput{}, nil

}

func shouldAddEntry(roll int) bool {
	return roll == 1 || roll == 6
}
