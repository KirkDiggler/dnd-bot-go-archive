package session

import (
	"context"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/session"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"
)

type manager struct {
	repo          session.Repository
	charManager   characters.Manager
}

type Config struct {
	Repo        session.Repository
	CharManager characters.Manager
}

func New(cfg *Config) (Manager, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Repo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Repo")
	}

	if cfg.CharManager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharManager")
	}

	return &manager{
		repo:          cfg.Repo,
		charManager:   cfg.CharManager,
	}, nil
}

func (m *manager) Create(ctx context.Context, userID string, draftID string) (*entities.Session, error) {
	if userID == "" {
		return nil, dnderr.NewMissingParameterError("userID")
	}
	if draftID == "" {
		return nil, dnderr.NewMissingParameterError("draftID")
	}

	session := &entities.Session{
		UserID:    userID,
		DraftID:   draftID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return m.repo.Create(ctx, session)
}

func (m *manager) GetWithDraft(ctx context.Context, userID string) (*entities.Session, error) {
	if userID == "" {
		return nil, dnderr.NewMissingParameterError("userID")
	}

	session, err := m.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	draft, err := m.charManager.GetDraft(ctx, session.DraftID)
	if err != nil {
		return nil, err
	}

	session.Draft = draft
	return session, nil
}

func (m *manager) Update(ctx context.Context, session *entities.Session) (*entities.Session, error) {
	if session == nil {
		return nil, dnderr.NewMissingParameterError("session")
	}

	session.UpdatedAt = time.Now()
	return m.repo.Update(ctx, session)
}

func (m *manager) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return dnderr.NewMissingParameterError("userID")
	}

	return m.repo.Delete(ctx, userID)
}

func (m *manager) UpdateLastToken(ctx context.Context, userID string, lastToken string) error {
	if userID == "" {
		return dnderr.NewMissingParameterError("userID")
	}
	if lastToken == "" {
		return dnderr.NewMissingParameterError("lastToken")
	}

	session, err := m.repo.Get(ctx, userID)
	if err != nil {
		return err
	}

	session.LastToken = lastToken
	session.UpdatedAt = time.Now()

	_, err = m.repo.Update(ctx, session)
	return err
}
