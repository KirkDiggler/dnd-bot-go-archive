package encounter

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client    redis.UniversalClient
	uuider    types.UUIDGenerator
	timeClock types.TimeClock
}

type RedisConfig struct {
	Client redis.UniversalClient
}

func NewRedis(cfg *RedisConfig) (*Redis, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &Redis{
		client:    cfg.Client,
		uuider:    &types.GoogleUUID{},
		timeClock: &types.Clock{},
	}, nil
}

func getEncounterKey(id string) string {
	return "encounter:" + id
}

func jsonToEncounter(jsonStr string) (*entities.Encounter, error) {
	if jsonStr == "" {
		return nil, dnderr.NewMissingParameterError("jsonStr")
	}

	var encounter entities.Encounter
	err := json.Unmarshal([]byte(jsonStr), &encounter)
	if err != nil {
		return nil, err
	}

	return &encounter, nil
}

func (r *Redis) Create(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error) {
	if encounter == nil {
		return nil, dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID != "" {
		return nil, dnderr.NewInvalidEntityError("encounter.ID")
	}

	encounter.ID = r.uuider.New()

	jsonStr, err := encounter.MarshallJSON()
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, getEncounterKey(encounter.ID), jsonStr, 0).Err()

	if err != nil {
		return nil, err
	}

	return encounter, nil
}

func (r *Redis) Update(ctx context.Context, encounter *entities.Encounter) (*entities.Encounter, error) {
	if encounter == nil {
		return nil, dnderr.NewMissingParameterError("encounter")
	}

	if encounter.ID == "" {
		return nil, dnderr.NewInvalidEntityError("encounter.ID cannot be null")
	}

	jsonStr, err := encounter.MarshallJSON()
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, getEncounterKey(encounter.ID), jsonStr, 0).Err()
	if err != nil {
		return nil, err
	}

	return encounter, nil
}

func (r *Redis) Get(ctx context.Context, id string) (*entities.Encounter, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	jsonStr, err := r.client.Get(ctx, getEncounterKey(id)).Result()
	if err != nil {
		return nil, err
	}

	encounter, err := jsonToEncounter(jsonStr)
	if err != nil {
		return nil, err
	}

	return encounter, nil
}
