package character

import (
	"context"
	"fmt"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/types"
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

func (r *redisRepo) GetCharacter(ctx context.Context, id string) (*Data, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result := r.client.Get(ctx, id)
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

	character.ID = r.uuider.New()

	result := r.client.Set(ctx, character.ID, dataToJSON(character), 0)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return character, nil
}
