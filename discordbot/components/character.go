package components

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"

	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/bwmarrin/discordgo"
)

const (
	selectCaracterAction    = "select-character"
	selectProficiencyAction = "select-proficiency"
	rollCharacterAction     = "roll-character"
	selectAttributeKey      = "select-attribute"
	buttonAttributeKey      = "button-attribute"
)

type Character struct {
	client      dnd5e.Client
	charManager characters.Manager
}

type CharacterConfig struct {
	Client        dnd5e.Client
	CharacterRepo characters.Manager
}

type charChoice struct {
	Name  string
	Race  *entities.Race
	Class *entities.Class
}

func NewCharacter(cfg *CharacterConfig) (*Character, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.CharacterRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharacterRepo")
	}
	return &Character{
		client:      cfg.Client,
		charManager: cfg.CharacterRepo,
	}, nil
}

func (c *Character) GetApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "character",
		Description: "Generate a character",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "random",
				Description: "Put a new character from a random list",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "load",
				Description: "Load your character",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c *Character) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "character":
			switch i.ApplicationCommandData().Options[0].Name {
			case "random":
				c.handleRandomStart(s, i)
			case "load":
				c.handleLoadCharacter(s, i)
			}
		}
	case discordgo.InteractionMessageComponent:
		strKey := fmt.Sprintf("%s:%s:str", selectAttributeKey, i.Member.User.ID)
		dexKey := fmt.Sprintf("%s:%s:dex", selectAttributeKey, i.Member.User.ID)
		conKey := fmt.Sprintf("%s:%s:con", selectAttributeKey, i.Member.User.ID)
		intKey := fmt.Sprintf("%s:%s:int", selectAttributeKey, i.Member.User.ID)
		wisKey := fmt.Sprintf("%s:%s:wis", selectAttributeKey, i.Member.User.ID)
		chaKey := fmt.Sprintf("%s:%s:cha", selectAttributeKey, i.Member.User.ID)

		switch i.MessageComponentData().CustomID {
		case selectCaracterAction:
			c.handleCharSelect(s, i)
		case rollCharacterAction:
			c.handleRollCharacter(s, i)
		case strKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "str", selectSlice)
		case dexKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "dex", selectSlice)
		case conKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "con", selectSlice)
		case intKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "int", selectSlice)
		case wisKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "wis", selectSlice)
		case chaKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "cha", selectSlice)
		case selectProficiencyAction:
			c.handleProficiencySelect(s, i)
		}
	}
}

// Setting Attributes
func (c *Character) handleAttributeSelect(s *discordgo.Session, i *discordgo.InteractionCreate, attribute string, selectSlice []string) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}
	idx, err := strconv.Atoi(selectSlice[1])
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	if idx >= len(char.Rolls) {
		log.Printf("index out of rabge. idx: %d, len: %d", idx, len(char.Rolls))
		return // TODO: Handle error
	}
	// TODO: make set attribut function that returns bool if it was set
	if !char.Rolls[idx].Used { // We have not used this one
		char.Attribues[entities.Attribute(attribute)].Score = char.Rolls[idx].Total - char.Rolls[idx].Lowest
		log.Println("setting ", attribute, " to ", char.Attribues[entities.Attribute(attribute)].Score)
		char.Rolls[idx].Used = true
		// TODO Calculate modifiers
	}

	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	rolls := char.Rolls
	attributeSelectData, err := c.generateAttributeSelect(char, rolls, i)
	done := false
	if err != nil {
		if err.Error() == "done" {
			done = true
		} else {
			log.Println(err)
			return // TODO: Handle error
		}
	}
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("Rolls: ")
	for _, roll := range rolls {
		msgBuilder.WriteString(fmt.Sprintf("%d, ", roll.Total-roll.Lowest))
	}

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	oldInteraction := &discordgo.Interaction{AppID: i.AppID, Token: state.LastToken}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	var response *discordgo.InteractionResponse
	if done {
		c.handleProficiencyStep(s, i)
		return
	} else {
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:      discordgo.MessageFlagsEphemeral,
				Content:    msgBuilder.String(),
				Components: attributeSelectData,
			},
		}
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
}

func (c *Character) generateAttributeSelect(char *entities.Character, rolls []*dice.RollResult, i *discordgo.InteractionCreate) ([]discordgo.MessageComponent, error) {
	userID := i.Member.User.ID

	selectionOrder := []string{"str", "dex", "con", "int", "wis", "cha"}

	selected := make(map[entities.Attribute]*entities.AbilityScore)

	for k, v := range char.Attribues {
		selected[k] = v
	}

	attrToSelect := ""
	for _, attr := range selectionOrder {
		if selected[entities.Attribute(attr)].Score == 0 {
			attrToSelect = attr
			break
		}
	}

	if attrToSelect == "" {
		return nil, errors.New("done")
	}

	components := make([]discordgo.SelectMenuOption, 0)
	for idx, roll := range rolls {
		if !roll.Used {
			components = append(components, discordgo.SelectMenuOption{
				Label: fmt.Sprintf("%v  %d", roll.Rolls, roll.Total-roll.Lowest),
				Value: fmt.Sprintf("roll:%d:%d", idx, roll.Total-roll.Lowest),
			})
		}
	}

	response := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "STR", selected[entities.AttributeStrength].Score),
					Disabled: selected[entities.AttributeStrength].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:str", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "DEX", selected[entities.AttributeDexterity].Score),
					Disabled: selected[entities.AttributeDexterity].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:dex", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "CON", selected[entities.AttributeConstitution].Score),
					Disabled: selected[entities.AttributeConstitution].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:con", buttonAttributeKey, userID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "INT", selected[entities.AttributeIntelligence].Score),
					Disabled: selected[entities.AttributeIntelligence].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:int", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "WIS", selected[entities.AttributeWisdom].Score),
					Disabled: selected[entities.AttributeWisdom].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:wis", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "CHA", selected[entities.AttributeCharisma].Score),
					Disabled: selected[entities.AttributeCharisma].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:cha", buttonAttributeKey, userID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					Placeholder: strings.ToUpper(attrToSelect),
					CustomID:    fmt.Sprintf("%s:%s:%s", selectAttributeKey, userID, attrToSelect),
					Options:     components,
				},
			},
		},
	}

	return response, nil
}

func (c *Character) handleRollCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	log.Println("Rolling for", i.Member.User.Username, "the ", char.Race.Name, " ", char.Class.Name)

	_, err = c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	rolls, err := rollAttributes()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	char.Rolls = rolls
	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	attributeSelectData, err := c.generateAttributeSelect(char, rolls, i)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("Rolls: ")
	for _, roll := range rolls {
		msgBuilder.WriteString(fmt.Sprintf("%d, ", roll.Total-roll.Lowest))
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Content:    msgBuilder.String(),
			Components: attributeSelectData,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
}

// Selecting a character
func (c *Character) handleRandomStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("handleRandomStart called")
	choices, err := c.startNewChoices(4)
	if err != nil {
		log.Println(err)
		// TODO: Handle error
		return
	}

	_, err = c.charManager.SaveState(context.Background(), &entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepSelect,
	})
	if err != nil {
		log.Println(err)
	}

	components := make([]discordgo.SelectMenuOption, len(choices))
	for idx, char := range choices {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%s the %s %s", i.Member.User.Username, char.Race.Name, char.Class.Name),
			Value: fmt.Sprintf("choice:%d:%s:%s", idx, char.Race.Key, char.Class.Key),
		}
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Select your new character:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    selectCaracterAction,
							Placeholder: "This is the tale of...👇",
							Options:     components,
						},
					},
				},
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}

}

func (c *Character) handleCharSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectString := strings.Split(i.MessageComponentData().Values[0], ":")
	if len(selectString) != 4 {
		log.Printf("Invalid select string: %s", selectString)
		return
	}

	race := selectString[2]
	class := selectString[3]

	_, err := c.charManager.Put(context.Background(), &entities.Character{
		OwnerID: i.Member.User.ID,
		Name:    i.Member.User.Username,
		Race: &entities.Race{
			Key: race,
		},
		Class: &entities.Class{
			Key: class,
		},
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	lastState, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	oldInteraction := &discordgo.Interaction{
		AppID: i.AppID,
		Token: lastState.LastToken,
	}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	c.handleRollCharacter(s, i)
}

// Selecting proficiency options
func (c *Character) generateProficiencyChoices(char *entities.Character, choices []*entities.Choice) (string, []discordgo.MessageComponent, error) {
	if char.Class == nil {
		log.Println("Class is nil")
		return "", nil, errors.New("class is nil")
	}

	if len(choices) == 0 {
		log.Println("No proficiency choices")
		return "", nil, errors.New("no proficiency choices")
	}

	var selectedChoice *entities.Choice
	for _, choice := range choices {
		if len(choice.Options) == 0 {
			log.Println("No proficiency choices")
			return "", nil, errors.New("no proficiency choices")
		}

		if choice.Status != entities.ChoiceStatusSelected {
			choice.Status = entities.ChoiceStatusActive
			selectedChoice = choice
			break
		}
	}

	err := c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency, choices)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	if selectedChoice == nil {
		log.Println("No proficiency choices")
		return "", nil, errors.New("no proficiency choices")
	}

	msg := fmt.Sprintf("Select %d starting proficiencies:", selectedChoice.Count)

	options := make([]discordgo.SelectMenuOption, len(selectedChoice.Options))
	for idx, choice := range selectedChoice.Options {
		options[idx] = discordgo.SelectMenuOption{
			Label: choice.GetName(),
			Value: fmt.Sprintf("choice:%s", choice.GetKey()),
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.SelectMenu{
			MinValues: &selectedChoice.Count,
			MaxValues: selectedChoice.Count,
			CustomID:  selectProficiencyAction,
			Options:   options,
		},
	}

	return msg, components, nil
}

func (c *Character) handleProficiencySelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepProficiency,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	choices, err := c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var done bool
	for _, choice := range choices {
		if choice.Status == entities.ChoiceStatusActive {
			choice.Status = entities.ChoiceStatusSelected

			// here is where I would reload the options for a given choice option
			for _, value := range i.MessageComponentData().Values {
				char.AddProficiency(&entities.Proficiency{
					Key: strings.Split(value, ":")[1],
				})
			}
			break
		}

		done = true
	}

	err = c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency, choices)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	oldInteraction := &discordgo.Interaction{
		AppID: i.AppID,
		Token: state.LastToken,
	}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if done {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Proficiencies selected",
			},
		}

		err = s.InteractionRespond(i.Interaction, response)
		if err != nil {
			log.Println(err)
			return // TODO handle error
		}
	} else {
		c.handleProficiencyStep(s, i)
	}
}

func (c *Character) handleProficiencyStep(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepProficiency,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var choices []*entities.Choice

	if state.Step == entities.CreateStepRoll {
		choices = char.Class.ProficiencyChoices

	} else {
		var choicesErr error
		choices, choicesErr = c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency)
		if choicesErr != nil {
			var notFoundErr *dnderr.NotFoundError
			if !errors.Is(choicesErr, notFoundErr) {
				log.Println(choicesErr)
				return // TODO handle error
			}
		}
	}

	msg, components, err := c.generateProficiencyChoices(char, choices)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: msg,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

// Gets the current state for returning before setting the input state
func (c *Character) getAndUpdateState(input *entities.CharacterCreation) (*entities.CharacterCreation, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	existing, err := c.charManager.GetState(context.Background(), input.CharacterID)
	if err != nil {
		return nil, err
	}

	_, err = c.charManager.SaveState(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (c *Character) handleLoadCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if char.Race == nil || char.Class == nil {
		log.Println("Character not fully loaded")
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Loaded character %s the %s %s", char.Name, char.Race.Name, char.Class.Name),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return
	}

}

func (c *Character) startNewChoices(number int) ([]*charChoice, error) {
	// TODO cache these. cache in the client wrapper? configurable?
	races, err := c.client.ListRaces()
	if err != nil {
		return nil, err
	}

	classes, err := c.client.ListClasses()
	if err != nil {
		return nil, err
	}
	log.Println("Starting new choices")
	choices := make([]*charChoice, number)

	var bard *entities.Class
	for _, class := range classes {
		if class.Key == "bard" {
			bard = class
			break
		}
	}
	for idx := 0; idx < number-1; idx++ {
		rand.Seed(time.Now().UnixNano())
		class := classes[rand.Intn(len(classes))]
		time.Sleep(1 * time.Nanosecond)
		rand.Seed(time.Now().UnixNano())
		race := races[rand.Intn(len(races))]
		choices[idx] = &charChoice{
			Race:  race,
			Class: class,
		}
	}
	choices[number-1] = &charChoice{
		Race:  races[rand.Intn(len(races))],
		Class: bard,
	}

	return choices, nil
}

func rollAttributes() ([]*dice.RollResult, error) {
	attributes := make([]*dice.RollResult, 6)

	for idx := range attributes {
		roll, err := dice.RollString("4d6")
		if err != nil {
			return nil, err
		}
		attributes[idx] = roll
	}
	return attributes, nil
}
