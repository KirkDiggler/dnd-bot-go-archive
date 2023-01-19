package party

import (
	"context"
	"errors"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/entities"

	"github.com/KirkDiggler/dnd-bot-go/types"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/suite"
)

type partySuite struct {
	suite.Suite

	ctx        context.Context
	fixture    *redisRepo
	redisMock  redismock.ClientMock
	mockUuider *types.MockUUID
}

func (s *partySuite) SetupTest() {
	s.ctx = context.Background()
	client, redisMock := redismock.NewClientMock()
	s.redisMock = redisMock
	s.mockUuider = &types.MockUUID{}

	s.fixture = &redisRepo{
		client: client,
		uuider: s.mockUuider,
	}
}

func (s *partySuite) TestGetParty() {
	jsonString := `{"name":"Test Party","party_size":5, "token":"1234"}`

	s.redisMock.ExpectGet("token").SetVal(jsonString)

	result, err := s.fixture.GetParty(s.ctx, "token")
	s.NoError(err)
	s.NotNil(result)
	s.Equal("Test Party", result.Name)
	s.Equal(5, result.PartySize)
	s.Equal("1234", result.Token)
}

func (s *partySuite) TestGetPartyNotFound() {
	s.redisMock.ExpectGet("1234").SetErr(redis.Nil)

	result, err := s.fixture.GetParty(s.ctx, "1234")
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, dnderr.NewNotFoundError("token: 1234 not found").Error())
}

func (s *partySuite) TestGetPartyError() {
	s.redisMock.ExpectGet("token").SetErr(errors.New("test error"))

	result, err := s.fixture.GetParty(s.ctx, "token")
	s.Error(err)
	s.Nil(result)
	s.EqualError(err, errors.New("test error").Error())
}

func (s *partySuite) TestCreateParty() {
	s.redisMock.ExpectSet("party:12345", `{"party_size":5,"name":"Test Party","token":"12345"}`, 0).SetVal("OK")
	s.mockUuider.On("New").Return("12345")

	result, err := s.fixture.CreateParty(s.ctx, &entities.Party{
		PartySize: 5,
		Name:      "Test Party",
	})

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Test Party", result.Name)
	s.Equal(5, result.PartySize)
	s.Equal("12345", result.Token)
}

func (s *partySuite) TestCreatePartyError() {
	s.redisMock.ExpectSet("party:12345", `{"party_size":5,"name":"Test Party","token":"12345"}`, 0).SetErr(errors.New("test error"))
	s.mockUuider.On("New").Return("12345")

	result, err := s.fixture.CreateParty(s.ctx, &entities.Party{
		PartySize: 5,
		Name:      "Test Party",
	})

	s.Error(err)
	s.Nil(result)

	s.EqualError(err, errors.New("test error").Error())
}
func TestParty(t *testing.T) {
	suite.Run(t, new(partySuite))
}
