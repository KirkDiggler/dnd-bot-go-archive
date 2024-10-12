package characters

import (
	"context"
	"log"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/encounter"
	"golang.org/x/sync/errgroup"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/choice"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_creation"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_draft"
)

type manager struct {
	charRepo      character.Repository
	stateRepo     character_creation.Repository
	choiceRepo    choice.Repository
	encounterRepo encounter.Repository
	client        dnd5e.Client
	draftRepo     character_draft.Repository
}

type Config struct {
	CharacterRepo character.Repository
	StateRepo     character_creation.Repository
	ChoiceRepo    choice.Repository
	Client        dnd5e.Client
	EncounterRepo encounter.Repository
	DraftRepo     character_draft.Repository
}

func New(cfg *Config) (Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.CharacterRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharacterRepo")
	}

	if cfg.StateRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.StateRepo")
	}

	if cfg.ChoiceRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.ChoiceRepo")
	}

	if cfg.EncounterRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.EncounterRepo")
	}

	if cfg.DraftRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.DraftRepo")
	}

	return &manager{
		charRepo:      cfg.CharacterRepo,
		stateRepo:     cfg.StateRepo,
		choiceRepo:    cfg.ChoiceRepo,
		client:        cfg.Client,
		encounterRepo: cfg.EncounterRepo,
		draftRepo:     cfg.DraftRepo,
	}, nil
}

func (m *manager) CreateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error) {
	if encounter == nil {
		return nil, dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID != "" {
		return nil, dnderr.NewInvalidParameterError("encounter.ID", encounter.ID)
	}

	result, err := m.encounterRepo.Create(ctx, encounter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) GetEncounter(ctx context.Context, id string) (*entities.Encounter, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result, err := m.encounterRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) UpdateEncounter(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error) {
	if encounter == nil {
		return nil, dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID == "" {
		return nil, dnderr.NewMissingParameterError("encounter.ID")
	}

	result, err := m.encounterRepo.Update(ctx, encounter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) AddProficiency(ctx context.Context, char *entities.Character, reference *entities.ReferenceItem) (*entities.Character, error) {
	if char == nil {
		return nil, dnderr.NewMissingParameterError("char")
	}

	if reference == nil {
		return nil, dnderr.NewMissingParameterError("reference")
	}

	if reference.Type != entities.ReferenceTypeProficiency {
		return nil, dnderr.NewInvalidParameterError("reference.Type", "must be a proficiency")
	}

	proficiency, err := m.client.GetProficiency(reference.Key)
	if err != nil {
		return nil, err
	}

	char.AddProficiency(proficiency)

	return char, nil
}

func (m *manager) AddInventory(ctx context.Context, char *entities.Character, key string) (*entities.Character, error) {
	if char == nil {
		return nil, dnderr.NewMissingParameterError("char")
	}

	if key == "" {
		return nil, dnderr.NewMissingParameterError("key")
	}

	equipment, err := m.client.GetEquipment(key)
	if err != nil {
		return nil, err
	}

	char.AddInventory(equipment)

	return char, nil
}

func (m *manager) SaveChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType, choices []*entities.Choice) error {
	if characterID == "" {
		return dnderr.NewMissingParameterError("characterID")
	}

	if choiceType == "" {
		return dnderr.NewMissingParameterError("choiceType")
	}

	if choices == nil {
		return dnderr.NewMissingParameterError("choices")
	}

	if len(choices) == 0 {
		return dnderr.NewMissingParameterError("choices")
	}

	return m.choiceRepo.Put(ctx, &choice.PutInput{
		CharacterID: characterID,
		Type:        choiceType,
		Choices:     choices,
	})
}

func (m *manager) GetChoices(ctx context.Context, characterID string, choiceType entities.ChoiceType) ([]*entities.Choice, error) {
	if characterID == "" {
		return nil, dnderr.NewMissingParameterError("characterID")
	}

	if choiceType == "" {
		return nil, dnderr.NewMissingParameterError("choiceType")
	}

	data, err := m.choiceRepo.Get(ctx, &choice.GetInput{
		CharacterID: characterID,
		Type:        choiceType,
	})
	if err != nil {
		return nil, err
	}

	return data.Choices, nil
}

func (m *manager) List(ctx context.Context, ownerID string, status ...entities.CharacterStatus) ([]*entities.Character, error) {
	if ownerID == "" {
		return nil, dnderr.NewMissingParameterError("ownerID")
	}

	charDatas, err := m.charRepo.ListByOwnerAndStatus(ctx, ownerID, status...)
	if err != nil {
		return nil, err
	}

	characters := make([]*entities.Character, len(charDatas))

	// use characterFromData in gou routines to hydrate the characters"
	g, gCtx := errgroup.WithContext(ctx)
	for i := range charDatas {
		i := i
		g.Go(func() error {
			char, err := m.characterFromData(gCtx, charDatas[i])
			if err != nil {
				return err
			}
			characters[i] = char
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return characters, nil
}

func (m *manager) Put(ctx context.Context, character *entities.Character) (*entities.Character, error) {
	if character == nil {
		return nil, dnderr.NewMissingParameterError("character")
	}
	//
	//if character.Name == "" {
	//	return nil, dnderr.NewMissingParameterError("character.Name")
	//}

	if character.OwnerID == "" {
		return nil, dnderr.NewMissingParameterError("character.OwnerID")
	}

	//if character.Race == nil {
	//	return nil, dnderr.NewMissingParameterError("character.Race")
	//}
	//
	//if character.Class == nil {
	//	return nil, dnderr.NewMissingParameterError("character.Class")
	//}

	if character.Race != nil {
		log.Println("character race: ", character.Race.Key)
	}

	if character.Class != nil {
		log.Println("character class: ", character.Class.Key)
	}

	data, err := m.charRepo.Put(ctx, character)
	if err != nil {
		return nil, err
	}

	character.ID = data.ID

	return character, nil
}

func (m *manager) Get(ctx context.Context, id string) (*entities.Character, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	data, err := m.charRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return m.characterFromData(ctx, data)
}

// TODO: Move to state manager
func (m *manager) SaveState(ctx context.Context, state *entities.CharacterCreation) (*entities.CharacterCreation, error) {
	if state == nil {
		return nil, dnderr.NewMissingParameterError("state")
	}

	if state.CharacterID == "" {
		return nil, dnderr.NewMissingParameterError("state.PlayerID")
	}

	if state.OwnerID == "" {
		return nil, dnderr.NewMissingParameterError("state.OwnerID")
	}

	result, err := m.stateRepo.Put(ctx, state)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) GetState(ctx context.Context, characterID string) (*entities.CharacterCreation, error) {
	if characterID == "" {
		return nil, dnderr.NewMissingParameterError("characterID")
	}

	result, err := m.stateRepo.Get(ctx, characterID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) CreateDraft(ctx context.Context, ownerID string) (*entities.CharacterDraft, error) {
	character := &entities.Character{
		OwnerID: ownerID,
		Status:  entities.CharacterStatusDraft,
	}

	savedChar, err := m.charRepo.Put(ctx, character)
	if err != nil {
		return nil, err
	}

	draft := &entities.CharacterDraft{
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CurrentStep: entities.SelectRaceStep,
		Character:   savedChar,
	}

	result, err := m.draftRepo.Create(ctx, draft)
	if err != nil {
		// If there's an error, we should delete the character we just created
		_ = m.charRepo.Delete(ctx, savedChar.ID)
		return nil, err
	}

	return result, nil
}

func (m *manager) GetDraft(ctx context.Context, draftID string) (*entities.CharacterDraft, error) {
	if draftID == "" {
		return nil, dnderr.NewMissingParameterError("draftID")
	}

	draft, err := m.draftRepo.Get(ctx, draftID)
	if err != nil {
		return nil, err
	}

	if draft.Character == nil {
		return nil, dnderr.NewInvalidEntityError("Character is nil in draft")
	}

	hydratedCharacter, err := m.characterFromData(ctx, draft.Character)
	if err != nil {
		return nil, err
	}

	draft.Character = hydratedCharacter
	return draft, nil
}

func (m *manager) UpdateDraft(ctx context.Context, draft *entities.CharacterDraft) (*entities.CharacterDraft, error) {
	if draft == nil {
		return nil, dnderr.NewMissingParameterError("draft")
	}

	if draft.ID == "" {
		return nil, dnderr.NewMissingParameterError("draft.ID")
	}

	if draft.Character == nil {
		return nil, dnderr.NewInvalidEntityError("Character is nil in draft")
	}

	draft.UpdatedAt = time.Now()

	// Save the character
	_, err := m.charRepo.Put(ctx, draft.Character)
	if err != nil {
		return nil, err
	}

	// Update the draft
	result, err := m.draftRepo.Update(ctx, draft)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *manager) ListDrafts(ctx context.Context, ownerID string) ([]*entities.CharacterDraft, error) {
	if ownerID == "" {
		return nil, dnderr.NewMissingParameterError("ownerID")
	}

	drafts, err := m.draftRepo.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return drafts, nil
}

func (m *manager) DeleteDraft(ctx context.Context, draftID string) error {
	if draftID == "" {
		return dnderr.NewMissingParameterError("draftID")
	}

	draft, err := m.draftRepo.Get(ctx, draftID)
	if err != nil {
		return err
	}

	if draft.Character == nil || draft.Character.ID == "" {
		return dnderr.NewInvalidEntityError("Character is nil or has no ID in draft")
	}

	// Delete the associated character
	err = m.charRepo.Delete(ctx, draft.Character.ID)
	if err != nil {
		return err
	}

	// Delete the draft
	return m.draftRepo.Delete(ctx, draftID)
}

func (m *manager) ActivateCharacter(ctx context.Context, draftID string) error {
	draft, err := m.GetDraft(ctx, draftID)
	if err != nil {
		return err
	}

	if !draft.AllStepsCompleted() {
		return dnderr.NewInvalidOperationError("Cannot activate character: creation not complete")
	}

	draft.Character.Status = entities.CharacterStatusActive
	_, err = m.charRepo.Put(ctx, draft.Character)
	if err != nil {
		return err
	}

	// Delete the draft after activating the character
	return m.DeleteDraft(ctx, draftID)
}

func contains(slice []entities.CharacterStatus, value entities.CharacterStatus) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Add a new function to complete a step
func (m *manager) CompleteStep(ctx context.Context, draftID string, step entities.CreateStep) error {
	draft, err := m.GetDraft(ctx, draftID)
	if err != nil {
		return err
	}

	draft.CompleteStep(step)
	draft.UpdatedAt = time.Now()

	_, err = m.draftRepo.Update(ctx, draft)
	return err
}

// Add a new function to check if all steps are completed
func (m *manager) IsCharacterCreationComplete(ctx context.Context, draftID string) (bool, error) {
	draft, err := m.GetDraft(ctx, draftID)
	if err != nil {
		return false, err
	}

	return draft.AllStepsCompleted(), nil
}
