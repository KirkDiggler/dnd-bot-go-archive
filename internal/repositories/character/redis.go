package character

import (
	"context"
	"encoding/json"
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

func getStatusCharactersKey(status entities.CharacterStatus) string {
	return fmt.Sprintf("status:%s:characters", status)
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

	if character.ID == "" {
		character.ID = r.uuider.New()
	}

	data, err := json.Marshal(character)
	if err != nil {
		return nil, err
	}

	pipe := r.client.Pipeline()
	pipe.Set(ctx, getCharacterKey(character.ID), data, 0)
	pipe.SAdd(ctx, getOwnerCharactersKey(character.OwnerID), character.ID)
	pipe.SAdd(ctx, getStatusCharactersKey(character.Status), character.ID)

	_, err = pipe.Exec(ctx)
	if err != nil {
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

func (r *redisRepo) Delete(ctx context.Context, id string) error {
	if id == "" {
		return dnderr.NewMissingParameterError("id")
	}

	// Get the character to retrieve owner and status
	char, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	// Remove from all relevant sets
	pipe := r.client.Pipeline()
	pipe.Del(ctx, getCharacterKey(id))
	pipe.SRem(ctx, getOwnerCharactersKey(char.OwnerID), id)
	pipe.SRem(ctx, getStatusCharactersKey(char.Status), id)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *redisRepo) ListByOwnerAndStatus(ctx context.Context, ownerID string, status ...entities.CharacterStatus) ([]*Data, error) {
	if ownerID == "" {
		return nil, dnderr.NewMissingParameterError("ownerID")
	}

	var characterIDs []string
	var err error

	if len(status) == 0 {
		// If no status is specified, get all characters for the owner
		characterIDs, err = r.client.SMembers(ctx, getOwnerCharactersKey(ownerID)).Result()
	} else {
		// If status is specified, get the intersection of owner's characters and characters with the specified status
		keys := make([]string, len(status)+1)
		keys[0] = getOwnerCharactersKey(ownerID)
		for i, s := range status {
			keys[i+1] = getStatusCharactersKey(s)
		}
		characterIDs, err = r.client.SInter(ctx, keys...).Result()
	}

	if err != nil {
		return nil, err
	}

	characters := make([]*Data, 0, len(characterIDs))
	for _, id := range characterIDs {
		char, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		characters = append(characters, char)
	}

	return characters, nil
}
