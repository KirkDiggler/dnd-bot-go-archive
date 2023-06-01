package encounter

import "context"

type Repository interface {
	CreateEncounter(ctx context.Context, encounter *Data) (*Data, error)
	UpdateEncounter(ctx context.Context, encounter *Data) (*Data, error)
	GetEncounter(ctx context.Context, id string) (*Data, error)
	ListEncountersByCharacter(ctx context.Context, owner string) ([]*Data, error)
}
