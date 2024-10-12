package characters

import (
	"context"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Manager interface {
	// Character operations
	Put(ctx context.Context, character *entities.Character) (*entities.Character, error)
	Get(ctx context.Context, id string) (*entities.Character, error)
	List(ctx context.Context, ownerID string, status ...entities.CharacterStatus) ([]*entities.Character, error)
	AddProficiency(ctx context.Context, char *entities.Character, reference *entities.ReferenceItem) (*entities.Character, error)
	AddInventory(ctx context.Context, char *entities.Character, key string) (*entities.Character, error)

	// Choice operations
	GetChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType) ([]*entities.Choice, error)
	SaveChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType, choices []*entities.Choice) error

	// Encounter operations
	CreateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	UpdateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error)
	GetEncounter(ctx context.Context, id string) (*entities.Encounter, error)

	// Draft operations
	CreateDraft(ctx context.Context, ownerID string) (*entities.CharacterDraft, error)
	GetDraft(ctx context.Context, draftID string) (*entities.CharacterDraft, error)
	UpdateDraft(ctx context.Context, draft *entities.CharacterDraft) (*entities.CharacterDraft, error)
	ListDrafts(ctx context.Context, ownerID string) ([]*entities.CharacterDraft, error)
	DeleteDraft(ctx context.Context, draftID string) error
	ActivateCharacter(ctx context.Context, draftID string) error
	CompleteStep(ctx context.Context, draftID string, step entities.CreateStep) error
	IsCharacterCreationComplete(ctx context.Context, draftID string) (bool, error)
}
