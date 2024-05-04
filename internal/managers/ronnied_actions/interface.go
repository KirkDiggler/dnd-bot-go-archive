package ronnied_actions

import "context"

type Interface interface {
	JoinGame(ctx context.Context, input *JoinGameInput) (*JoinGameOutput, error)
	AddRoll(ctx context.Context, input *AddRollInput) (*AddRollOutput, error)
	AddRolls(ctx context.Context, input *AddRollsInput) (*AddRollsOutput, error)
	GetTab(ctx context.Context, input *GetTabInput) (*GetTabOutput, error)
	ListTabs(ctx context.Context, input *ListTabsInput) (*ListTabsOutput, error)
	PayDrink(ctx context.Context, input *PayDrinkInput) (*PayDrinkOutput, error)
	CreateSession(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error)
	JoinSession(ctx context.Context, input *JoinSessionInput) (*JoinSessionOutput, error)
	GetSession(ctx context.Context, input *GetSessionInput) (*GetSessionOutput, error)
}
