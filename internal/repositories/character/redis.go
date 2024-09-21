package character

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/internal/types"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
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

func getOwnerCharactersKey(ownerID string) string {
	return fmt.Sprintf("owner:%s:characters", ownerID)
}

func (r *redisRepo) Get(ctx context.Context, id string) (*Data, error) {
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

	data := jsonToData(result.Val())

	return data, nil
}

func (r *redisRepo) Put(ctx context.Context, character *entities.Character) (*entities.Character, error) {
	if character == nil {
		return nil, dnderr.NewMissingParameterError("character")
	}

	if character.OwnerID == "" {
		return nil, dnderr.NewMissingParameterError("character.OwnerID")
	}

	character.ID = r.uuider.New()
	data := dataToJSON(characterToData(character))

	result := r.client.Set(ctx, getCharacterKey(character.ID), data, 0)
	if result.Err() != nil {
		return nil, result.Err()
	}

	// Add character ID to the owner's list of characters
	ownerCharactersKey := getOwnerCharactersKey(character.OwnerID)
	if err := r.client.SAdd(ctx, ownerCharactersKey, character.ID).Err(); err != nil {
		return nil, err
	}

	return character, nil
}

func (r *redisRepo) ListByOwner(ctx context.Context, ownerID string) ([]*Data, error) {
	if ownerID == "" {
		return nil, dnderr.NewMissingParameterError("ownerID")
	}

	ownerCharactersKey := getOwnerCharactersKey(ownerID)
	characterIDs, err := r.client.SMembers(ctx, ownerCharactersKey).Result()
	if err != nil {
		return nil, err
	}

	var characters []*Data
	for _, id := range characterIDs {
		character, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}
