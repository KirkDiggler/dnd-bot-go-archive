package ronnied_actions

import (
	"context"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/ronnied/game"
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

func (m *Manager) JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	result, err := m.gameRepo.Join(ctx, &game.JoinInput{
		GameID:   input.GameID,
		MemberID: input.MemberID,
	})
	if err != nil {
		return nil, err
	}

	return &JoinGameOutput{
		Member: result.Member,
	}, nil
}

func (m *Manager) AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.GameID == "" {
		return nil, dnderr.NewMissingParameterError("input.GameID")
	}

	if input.MemberID == "" {
		return nil, dnderr.NewMissingParameterError("input.MemberID")
	}

	if shouldAddEntry(input.Roll) {
		_, err := m.gameRepo.AddEntry(ctx, &game.AddEntryInput{
			GameID:   input.GameID,
			MemberID: input.MemberID,
			Roll:     input.Roll,
		})
		if err != nil {
			return nil, err
		}
	}

	return &AddRollOutput{}, nil
}

func shouldAddEntry(roll int) bool {
	return roll == 1 || roll == 6
}
