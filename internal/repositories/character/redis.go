package character

import (
	"context"
	"fmt"

	"github.com/KirkDiggler/dnd-bot-go/internal/types"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/go-redis/redis/v9"
)

type redisRepo struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
}

type Config struct {
	Client redis.UniversalClient
}

func New(cfg *Config) (Repository, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &redisRepo{
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}
func getCharacterKey(id string) string {
	return fmt.Sprintf("character:%s", id)
}

func (r *redisRepo) GetCharacter(ctx context.Context, id string) (*Data, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result := r.client.Get(ctx, getCharacterKey(id))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError(fmt.Sprintf("character id not found: %s", id))
		}

		return nil, result.Err()
	}

	return jsonToData(result.Val()), nil
}

func (r *redisRepo) CreateCharacter(ctx context.Context, character *Data) (*Data, error) {
	if character == nil {
		return nil, dnderr.NewMissingParameterError("character")
	}

	if character.OwnerID == "" {
		return nil, dnderr.NewMissingParameterError("character.OwnerID")
	}

	character.ID = character.OwnerID

	result := r.client.Set(ctx, getCharacterKey(character.ID), dataToJSON(character), 0)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return character, nil
}
