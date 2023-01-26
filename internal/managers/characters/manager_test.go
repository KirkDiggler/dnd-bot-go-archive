package characters

import (
	"context"
	"errors"
	"testing"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character_creation"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"

	"github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"
	"github.com/stretchr/testify/suite"
)

type managerSuite struct {
	suite.Suite

	ctx           context.Context
	fixture       *manager
	mockRepo      *character.Mock
	mockStateRepo *character_creation.Mock
	mockClient    *dnd5e.Mock
	id            string
	race          *entities.Race
	class         *entities.Class
	character     *entities.Character
	characterData *character.Data
}

func (s *managerSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockRepo = &character.Mock{}
	s.mockStateRepo = &character_creation.Mock{}
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
		Rolls: make([]*dice.RollResult, 0),
	}
	s.characterData = characterToData(s.character)
	s.fixture = &manager{
		charRepo:  s.mockRepo,
		client:    s.mockClient,
		stateRepo: s.mockStateRepo,
	}
}

func (s *managerSuite) TestSaveState() {
	state := &entities.CharacterCreation{
		CharacterID: s.id,
		LastToken:   "token",
		Step:        entities.CreateStepProficiency,
	}

	s.mockStateRepo.On("Put", s.ctx, state).Return(
		state, nil)

	_, err := s.fixture.SaveState(s.ctx, state)
	s.NoError(err)
}

func (s *managerSuite) TestSaveStateMissingState() {
	_, err := s.fixture.SaveState(s.ctx, nil)
	s.Error(err)
	s.EqualError(err, "Missing parameter: state")
}

func (s *managerSuite) TestSaveStateMissingCharacterID() {
	_, err := s.fixture.SaveState(s.ctx, &entities.CharacterCreation{})
	s.Error(err)
	s.EqualError(err, "Missing parameter: state.CharacterID")
}

func (s *managerSuite) TestSaveStateRepoErrors() {
	state := &entities.CharacterCreation{
		CharacterID: s.id,
		LastToken:   "token",
		Step:        entities.CreateStepProficiency,
	}

	s.mockStateRepo.On("Put", s.ctx, state).Return(
		nil, errors.New("test error"))

	_, err := s.fixture.SaveState(s.ctx, state)
	s.Error(err)
	s.EqualError(err, "test error")
}

func (s *managerSuite) TestGetState() {
	state := &entities.CharacterCreation{
		CharacterID: s.id,
		LastToken:   "token",
		Step:        entities.CreateStepProficiency,
	}

	s.mockStateRepo.On("Get", s.ctx, s.id).Return(
		state, nil)

	actual, err := s.fixture.GetState(s.ctx, s.id)
	s.NoError(err)
	s.Equal(state, actual)
}

func (s *managerSuite) TestGetStateMissingCharacterID() {
	_, err := s.fixture.GetState(s.ctx, "")
	s.Error(err)
	s.EqualError(err, "Missing parameter: characterID")
}

func (s *managerSuite) TestGetStateRepoErrors() {
	s.mockStateRepo.On("Get", s.ctx, s.id).Return(
		nil, errors.New("test error"))

	_, err := s.fixture.GetState(s.ctx, s.id)
	s.Error(err)
	s.EqualError(err, "test error")
}

func (s *managerSuite) TestCreate() {
	s.mockRepo.On("Put", s.ctx, s.character).Return(
		s.character, nil)

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
	s.mockRepo.On("Put", s.ctx, s.character).Return(
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

	s.mockRepo.On("Get", s.ctx, s.id).Return(s.characterData, nil)

	char, err := s.fixture.Get(s.ctx, s.id)
	s.NoError(err)
	s.Equal(s.character, char)
}

func TestCharacter(t *testing.T) {
	suite.Run(t, new(managerSuite))
}
