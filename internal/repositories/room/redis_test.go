package room

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
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

	room     *Data
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
	s.room = &Data{
		ID:          "1234",
		Status:      RoomStatusActive,
		MonsterID:   "1234",
		CharacterID: "1337",
	}

	buf, _ := json.Marshal(s.room)
	s.roomJson = string(buf)
}

func (s *roomSuite) TestCreaeRoom_ValidateInput() {

	_, err := s.fixture.CreateRoom(s.ctx, nil)
	s.Error(err)
	s.EqualError(err, dnderr.NewMissingParameterError("room").Error())

	_, err = s.fixture.CreateRoom(s.ctx, &Data{ID: "1234"})
	s.Error(err)
	s.EqualError(err, dnderr.NewInvalidEntityError("room.ID must be empty").Error())
}

func (s *roomSuite) TestCreateRoom_RedisError() {
	s.mockUuider.On("New").Return(s.room.ID)

	s.redisMock.ExpectZCard(characterRoomKey(s.room.CharacterID)).SetErr(errors.New("redis error"))

	input := &Data{
		Status:      RoomStatusActive,
		MonsterID:   s.room.MonsterID,
		CharacterID: s.room.CharacterID,
	}
	result, err := s.fixture.CreateRoom(s.ctx, input)
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, "redis error")
}

func (s *roomSuite) TestCreateRoom() {
	s.mockUuider.On("New").Return(s.room.ID)

	s.redisMock.ExpectZCard(characterRoomKey(s.room.CharacterID)).SetVal(42)

	s.redisMock.ExpectTxPipeline()
	s.redisMock.ExpectSet(getRoomKey(s.room.ID), s.roomJson, 0).SetVal(s.roomJson)

	s.redisMock.ExpectZAdd(characterRoomKey(s.room.CharacterID), redis.Z{
		Score:  42,
		Member: getRoomKey(s.room.ID),
	}).SetVal(1)

	s.redisMock.ExpectTxPipelineExec()

	input := &Data{
		Status:      RoomStatusActive,
		MonsterID:   s.room.MonsterID,
		CharacterID: s.room.CharacterID,
	}
	result, err := s.fixture.CreateRoom(s.ctx, input)
	s.NoError(err)
	s.NotNil(result)
	s.Equal(s.room, result)
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
