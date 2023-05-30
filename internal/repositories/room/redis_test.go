package room

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/KirkDiggler/dnd-bot-go/internal/types"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/suite"
	"testing"
)

type roomSuite struct {
	suite.Suite

	ctx        context.Context
	redisMock  redismock.ClientMock
	mockUuider *types.MockUUID
	fixture    *Redis

	room     *entities.Room
	roomJson string
}

func (s *roomSuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.redisMock = redisMock
	s.mockUuider = &types.MockUUID{}

	s.fixture = &Redis{
		client: client,
		uuider: s.mockUuider,
	}
	s.room = &entities.Room{
		ID:          "1234",
		Status:      entities.RoomStatusActive,
		MonsterID:   "1234",
		CharacterID: "1337",
	}

	buf, _ := json.Marshal(s.room)
	s.roomJson = string(buf)
}

func (s *roomSuite) TestGetRoom() {
	s.redisMock.ExpectGet(getRoomKey(s.room.ID)).SetVal(s.roomJson)

	result, err := s.fixture.GetRoom(s.ctx, s.room.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.room, result)
}

func (s *roomSuite) TestGetRoomNotFound() {
	s.redisMock.ExpectGet(s.room.ID).SetErr(redis.Nil)

	result, err := s.fixture.GetRoom(s.ctx, s.room.ID)
	s.Error(err)
	s.Nil(result)
}

func (s *roomSuite) TestUpdateRoom() {
	s.redisMock.ExpectSet(getRoomKey(s.room.ID), s.roomJson, 0).SetVal(s.roomJson)

	result, err := s.fixture.UpdateRoom(s.ctx, s.room)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.room, result)
}

func (s *roomSuite) TestUpdateRoomError() {
	s.redisMock.ExpectSet(getRoomKey(s.room.ID), s.roomJson, 0).SetErr(errors.New("error"))

	result, err := s.fixture.UpdateRoom(s.ctx, s.room)
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, "error")
}

func TestRoom(t *testing.T) {
	suite.Run(t, new(roomSuite))
}
