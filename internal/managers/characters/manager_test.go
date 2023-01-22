package characters

import (
	"context"
	"errors"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
	"github.com/stretchr/testify/suite"
)

type managerSuite struct {
	suite.Suite

	ctx        context.Context
	fixture    *manager
	mockRepo   *character.Mock
	mockClient *dnd5e.Mock
	id         string
	race       *entities.Race
	class      *entities.Class
	character  *entities.Character
}

func (s *managerSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockRepo = &character.Mock{}
	s.mockClient = &dnd5e.Mock{}
	s.id = "123"
	s.race = &entities.Race{
		Key:  "elf",
		Name: "Elf",
	}
	s.class = &entities.Class{
		Key:  "fighter",
		Name: "Fighter",
	}
	s.character = &entities.Character{
		ID:      s.id,
		Name:    "Test Character",
		OwnerID: s.id,
		Race:    s.race,
		Class:   s.class,
		Attribues: map[entities.Attribute]*entities.AbilityScore{
			entities.AttributeStrength:     {Score: 16},
			entities.AttributeDexterity:    {Score: 15},
			entities.AttributeConstitution: {Score: 14},
			entities.AttributeIntelligence: {Score: 13},
			entities.AttributeWisdom:       {Score: 12},
			entities.AttributeCharisma:     {Score: 11},
		},
	}
	s.fixture = &manager{
		charRepo: s.mockRepo,
		client:   s.mockClient,
	}
}

func (s *managerSuite) TestCreate() {
	s.mockRepo.On("Put", s.ctx, s.character.ToData()).Return(
		s.character.ToData(), nil)

	char, err := s.fixture.Put(s.ctx, s.character)
	s.NoError(err)
	s.Equal(s.character, char)
}

func (s *managerSuite) TestCreateMissingCharacter() {
	_, err := s.fixture.Put(s.ctx, nil)
	s.Error(err)
	s.EqualError(err, "Missing parameter: character")
}

func (s *managerSuite) TestCreateMissingCharacterName() {
	s.character.Name = ""
	_, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "Missing parameter: character.Name")
}

func (s *managerSuite) TestCreateMissingCharacterOwnerID() {
	s.character.OwnerID = ""
	_, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "Missing parameter: character.OwnerID")
}

func (s *managerSuite) TestCreateMissingCharacterRace() {
	s.character.Race = nil
	_, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "Missing parameter: character.Race")
}

func (s *managerSuite) TestCreateMissingCharacterClass() {
	s.character.Class = nil
	_, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "Missing parameter: character.Class")
}

func (s *managerSuite) TestCreateRepoErrors() {
	s.mockRepo.On("Put", s.ctx, s.character.ToData()).Return(
		nil, errors.New("repo error"))

	_, err := s.fixture.Put(s.ctx, s.character)
	s.Error(err)
	s.EqualError(err, "repo error")
}

func (s *managerSuite) TestGet() {
	s.mockClient.On("GetRace", s.race.Key).Return(
		s.race, nil)
	s.mockClient.On("GetClass", s.class.Key).Return(
		s.class, nil)

	s.mockRepo.On("Get", s.ctx, s.id).Return(s.character.ToData(), nil)

	char, err := s.fixture.Get(s.ctx, s.id)
	s.NoError(err)
	s.Equal(s.character, char)
}

func TestCharacter(t *testing.T) {
	suite.Run(t, new(managerSuite))
}
