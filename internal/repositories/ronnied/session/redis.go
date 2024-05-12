package session

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Redis struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
}

type Config struct {
	Client redis.UniversalClient
}

func NewRedis(cfg *Config) (*Redis, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &Redis{
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}

func (r *Redis) Get(ctx context.Context, input *GetInput) (*GetOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.ID == "" {
		return nil, dnderr.NewMissingParameterError("input.ID")
	}

	sessionKey := getSessionKey(input.ID)

	sessionJson, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	session := &ronnied.Session{}
	err = json.Unmarshal([]byte(sessionJson), session)
	if err != nil {
		return nil, err
	}

	return &GetOutput{
		Session: session,
	}, nil
}

func (r *Redis) Create(ctx context.Context, input *CreateInput) (*CreateOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	now := time.Now()

	session := &ronnied.Session{
		ID:          r.uuider.New(),
		GameID:      input.GameID,
		SessionDate: &now,
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	sessionKey := getSessionKey(session.ID)

	err = r.client.Set(ctx, sessionKey, sessionJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return &CreateOutput{
		Session: session,
	}, nil
}

func (r *Redis) Update(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.Session == nil {
		return nil, dnderr.NewMissingParameterError("input.Session")
	}

	if input.Session.ID == "" {
		return nil, dnderr.NewMissingParameterError("input.Session.ID")
	}

	sessionJson, err := json.Marshal(input.Session)
	if err != nil {
		return nil, err
	}

	sessionKey := getSessionKey(input.Session.ID)

	err = r.client.Set(ctx, sessionKey, sessionJson, 0).Err()
	if err != nil {
		return nil, err
	}

	return &UpdateOutput{
		Session: input.Session,
	}, nil
}

func (r *Redis) JoinSessionRoll(ctx context.Context, input *JoinSessionRollInput) (*JoinSessionRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionRollID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionRollID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	rollKey := getSessionRollKey(input.SessionRollID)

	rollJson, err := r.client.Get(ctx, rollKey).Result()
	if err != nil {
		return nil, err
	}

	roll := &ronnied.SessionRoll{}
	err = json.Unmarshal([]byte(rollJson), roll)
	if err != nil {
		return nil, err
	}

	if roll.HasPlayer(input.PlayerID) {
		return nil, dnderr.NewAlreadyExistsError("player already in roll")
	}

	roll.Players = append(roll.Players, input.PlayerID)

	rollBytes, err := json.Marshal(roll)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, rollKey, string(rollBytes), 0).Err()
	if err != nil {
		return nil, err
	}

	return &JoinSessionRollOutput{
		SessionRoll: roll,
	}, nil
}

// Join adds a player to a session by setting them to the roll type start session roll
func (r *Redis) Join(ctx context.Context, input *JoinInput) (*JoinOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	sessionKey := getSessionKey(input.SessionID)

	sessionJson, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	session := &ronnied.Session{}
	err = json.Unmarshal([]byte(sessionJson), session)
	if err != nil {
		return nil, err
	}

	if session.HasPlayer(input.PlayerID) {
		return nil, dnderr.NewAlreadyExistsError("player already in session")
	}

	session.Players = append(session.Players, input.PlayerID)

	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, sessionKey, string(sessionBytes), 0).Err()
	if err != nil {
		return nil, err
	}

	return &JoinOutput{
		Session: session,
	}, nil
}

func (r *Redis) CreateRoll(ctx context.Context, input *CreateRollInput) (*CreateRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionID")
	}

	sessionKey := getSessionKey(input.SessionID)

	sessionJson, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	session := &ronnied.Session{}
	err = json.Unmarshal([]byte(sessionJson), session)
	if err != nil {
		return nil, err
	}

	roll := &ronnied.SessionRoll{
		ID:        r.uuider.New(),
		SessionID: session.ID,
		Type:      ronnied.RollTypeStart,
		Players:   session.Players,
	}

	rollBytes, err := json.Marshal(roll)
	if err != nil {
		return nil, err
	}

	rollKey := getSessionRollKey(roll.ID)

	err = r.client.Set(ctx, rollKey, string(rollBytes), 0).Err()
	if err != nil {
		return nil, err
	}

	return &CreateRollOutput{
		SessionRoll: roll,
	}, nil
}

func (r *Redis) GetSessionRoll(ctx context.Context, input *GetSessionRollInput) (*GetSessionRollOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.ID == "" {
		return nil, dnderr.NewMissingParameterError("input.ID")
	}

	rollKey := getSessionRollKey(input.ID)

	rollJson, err := r.client.Get(ctx, rollKey).Result()
	if err != nil {
		return nil, err
	}

	roll := &ronnied.SessionRoll{}
	err = json.Unmarshal([]byte(rollJson), roll)
	if err != nil {
		return nil, err
	}

	return &GetSessionRollOutput{
		SessionRoll: roll,
	}, nil
}

func (r *Redis) AddEntry(ctx context.Context, input *AddEntryInput) (*AddEntryOutput, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	if input.SessionRollID == "" {
		return nil, dnderr.NewMissingParameterError("input.SessionRollID")
	}

	if input.PlayerID == "" {
		return nil, dnderr.NewMissingParameterError("input.PlayerID")
	}

	sessionRollKey := getSessionRollKey(input.SessionRollID)

	sessionJson, err := r.client.Get(ctx, sessionRollKey).Result()
	if err != nil {
		log.Println("Failed to get session roll:", sessionRollKey, err)

		return nil, err
	}

	sessionRoll := &ronnied.SessionRoll{}
	err = json.Unmarshal([]byte(sessionJson), sessionRoll)
	if err != nil {
		return nil, err
	}

	if !sessionRoll.HasPlayer(input.PlayerID) {
		return nil, dnderr.NewNotFoundError("player not in session")
	}

	entry := &ronnied.SessionEntry{
		ID:            r.uuider.New(),
		SessionRollID: sessionRoll.ID,
		PlayerID:      input.PlayerID,
		Roll:          input.Roll,
		AssignedTo:    input.AssignedTo,
	}

	sessionRoll.Entries = append(sessionRoll.Entries, entry)

	entryBytes, err := json.Marshal(sessionRoll)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, sessionRollKey, string(entryBytes), 0).Err()
	if err != nil {
		return nil, err
	}

	return &AddEntryOutput{
		SessionEntry: entry,
	}, nil
}

func getSessionKey(id string) string {
	return "session:" + id
}
func getSessionRollKey(SessionRollID string) string {
	return "session:roll:" + SessionRollID
}
func getSessionEntryKey(SessionEntryID string) string {
	return "session:entry:" + SessionEntryID
}
