package ronnied_actions

import "context"

type Interface interface {
	JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error)
	AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error)
	GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error)
	PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error)
}
