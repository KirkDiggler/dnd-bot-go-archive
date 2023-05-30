package room

import (
	"context"
	"encoding/json"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/go-redis/redis/v9"
)

type Redis struct {
	client redis.UniversalClient
	uuider types.UUIDGenerator
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
		client: cfg.Client,
		uuider: &types.GoogleUUID{},
	}, nil
}

func characterRoomKey(characterID string) string {
	return "characterRoom:" + characterID
}

func getRoomKey(id string) string {
	return "room:" + id
}

func roomToJson(room *entities.Room) (string, error) {
	if room == nil {
		return "", dnderr.NewMissingParameterError("room")
	}

	buf, err := json.Marshal(room)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func jsonToRoom(jsonStr string) (*entities.Room, error) {
	if jsonStr == "" {
		return nil, dnderr.NewMissingParameterError("jsonStr")
	}

	var room entities.Room
	err := json.Unmarshal([]byte(jsonStr), &room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// CreateRoom creates a room and assigns it to the oweners index
func (r *Redis) CreateRoom(ctx context.Context, room *entities.Room) (*entities.Room, error) {
	if room == nil {
		return nil, dnderr.NewMissingParameterError("room")
	}

	if room.ID != "" {
		return nil, dnderr.NewInvalidEntityError("room.ID must be empty")
	}

	room.ID = r.uuider.New()

	jsonStr, err := roomToJson(room)
	if err != nil {
		return nil, err
	}

	roomCount, err := r.client.ZCard(ctx, characterRoomKey(room.CharacterID)).Result()
	if err != nil {
		return nil, err
	}

	// Use pipe to set the room data and add member to character rooms index
	pipe := r.client.TxPipeline()
	pipe.Set(ctx, getRoomKey(room.ID), jsonStr, 0)
	pipe.ZAdd(ctx, characterRoomKey(room.CharacterID), redis.Z{
		Score:  float64(roomCount),
		Member: getRoomKey(room.ID),
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Redis) UpdateRoom(ctx context.Context, room *entities.Room) (*entities.Room, error) {
	if room == nil {
		return nil, dnderr.NewMissingParameterError("room")
	}

	if room.ID == "" {
		return nil, dnderr.NewInvalidEntityError("room.ID must not be empty")
	}

	jsonStr, err := roomToJson(room)
	if err != nil {
		return nil, err
	}

	err = r.client.Set(ctx, getRoomKey(room.ID), jsonStr, 0).Err()
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *Redis) GetRoom(ctx context.Context, id string) (*entities.Room, error) {
	if id == "" {
		return nil, dnderr.NewMissingParameterError("id")
	}

	result := r.client.Get(ctx, getRoomKey(id))
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, dnderr.NewNotFoundError("room not found")
		}

		return nil, result.Err()
	}

	return jsonToRoom(result.Val())
}

func (r *Redis) ListRooms(ctx context.Context, owner string) ([]*entities.Room, error) {
	if owner == "" {
		return nil, dnderr.NewMissingParameterError("owner")
	}

	roomKeys, err := r.client.ZRevRange(ctx, characterRoomKey(owner), 0, 10).Result()
	if err != nil {
		return nil, err
	}

	rooms := make([]*entities.Room, len(roomKeys))
	for i, roomKey := range roomKeys {
		room, err := r.GetRoom(ctx, roomKey)
		if err != nil {
			return nil, err
		}

		rooms[i] = room
	}

	return rooms, nil
}
