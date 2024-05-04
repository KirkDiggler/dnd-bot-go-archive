package session

import "context"

type Interface interface {
	Create(ctx context.Context, input *CreateInput) (*CreateOutput, error)
	Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error)
	Get(ctx context.Context, input *GetInput) (*GetOutput, error)
	Join(ctx context.Context, input *JoinInput) (*JoinOutput, error)
	CreateRoll(ctx context.Context, input *CreateRollInput) (*CreateRollOutput, error)
	AddEntry(ctx context.Context, input *AddEntryInput) (*AddEntryOutput, error)
}
