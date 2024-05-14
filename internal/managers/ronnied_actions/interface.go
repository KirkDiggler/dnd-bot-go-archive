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
	UpdateSession(ctx context.Context, input *UpdateSessionInput) (*UpdateSessionOutput, error)
	JoinSession(ctx context.Context, input *JoinSessionInput) (*JoinSessionOutput, error)
	AddSessionRoll(ctx context.Context, input *AddSessionRollInput) (*AddSessionRollOutput, error)
	CreateSessionRoll(ctx context.Context, input *CreateSessionRollInput) (*CreateSessionRollOutput, error)
	UpdateSessionRoll(ctx context.Context, input *UpdateSessionRollInput) (*UpdateSessionRollOutput, error)
	GetSessionRoll(ctx context.Context, input *GetSessionRollInput) (*GetSessionRollOutput, error)
	GetSession(ctx context.Context, input *GetSessionInput) (*GetSessionOutput, error)
	AssignDrink(ctx context.Context, input *AssignDrinkInput) (*AssignDrinkOutput, error)
}
