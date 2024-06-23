package session

import "context"

type Interface interface {
	Create(ctx context.Context, input *CreateInput) (*CreateOutput, error)
	Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error)
	Get(ctx context.Context, input *GetInput) (*GetOutput, error)
	Join(ctx context.Context, input *JoinInput) (*JoinOutput, error)
	JoinSessionRoll(ctx context.Context, input *JoinSessionRollInput) (*JoinSessionRollOutput, error)
	CreateRoll(ctx context.Context, input *CreateRollInput) (*CreateRollOutput, error)
	UpdateRoll(ctx context.Context, input *UpdateRollInput) (*UpdateRollOutput, error)
	GetSessionRoll(ctx context.Context, input *GetSessionRollInput) (*GetSessionRollOutput, error)
	AddEntry(ctx context.Context, input *AddEntryInput) (*AddEntryOutput, error)
	ListSessionRolls(ctx context.Context, input *ListSessionRollsInput) (*ListSessionRollsOutput, error)
}
