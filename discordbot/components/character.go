package components

import (
	"context"
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
	selectCaracterAction        = "select-character"
	rollCharacterActionPhysical = "roll-character-physical"
	rollCharacterActionMental   = "roll-character-mental"
	selectAttributeKey          = "select-attribute"
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
		case rollCharacterActionPhysical:
			c.handleRollCharacterPhysical(s, i)
		case rollCharacterActionMental:
			c.handleRollCharacterMental(s, i)
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
		}
	}
}
func (c *Character) handleAttributeSelect(s *discordgo.Session, i *discordgo.InteractionCreate, attribute string, selectSlice []string) {
	log.Printf("User: %s, Attribute: %s, SelectSlice: %v", i.Member.User.ID, attribute, selectSlice)

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

	// TODO: make set attribut function that returns bool if it was set
	if !char.Rolls[idx].Used { // We have not used this one
		char.Attribues[entities.Attribute(attribute)].Score = char.Rolls[idx].Total - char.Rolls[idx].Lowest
		char.Rolls[idx].Used = true
		// TODO Calculate modifiers
	}

	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

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

func (c *Character) handleCharSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectString := strings.Split(i.MessageComponentData().Values[0], ":")
	if len(selectString) != 4 {
		log.Printf("Invalid select string: %s", selectString)
		return
	}

	race := selectString[2]
	class := selectString[3]

	char, err := c.charManager.Put(context.Background(), &entities.Character{
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

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Created character %s", char.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Roll Character",
							Style:    discordgo.SuccessButton,
							CustomID: rollCharacterActionPhysical,
							Emoji: discordgo.ComponentEmoji{
								Name: "🎲",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Character) GenerateAttributSelect(a *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	log.Println("Rolling character")
	rolls, err := rollAttributes()
	if err != nil {
		log.Println(err)
		return nil // TODO handle error
	}

	userID := i.Member.User.ID

	components := make([]discordgo.SelectMenuOption, len(rolls))
	for idx, roll := range rolls {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%v  %d", roll.Details, roll.Roll),
			Value: fmt.Sprintf("roll:%d:%d", roll.Index, roll.Roll),
		}
	}

	response := &discordgo.InteractionResponseData{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: "Assign your physical attributes:",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.SelectMenu{
						Placeholder: "STR",
						CustomID:    fmt.Sprintf("%s:%s:str", selectAttributeKey, userID),
						Options:     components,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.SelectMenu{
						Placeholder: "DEX",
						CustomID:    fmt.Sprintf("%s:%s:dex", selectAttributeKey, userID),
						Options:     components,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.SelectMenu{
						Placeholder: "CON",
						CustomID:    fmt.Sprintf("%s:%s:con", selectAttributeKey, userID),
						Options:     components,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Continue Non-Physical",
						Style:    discordgo.SuccessButton,
						CustomID: rollCharacterActionMental,
						Emoji: discordgo.ComponentEmoji{
							Name: "🎲",
						},
					},
				},
			},
		},
	}

	return response
}
func (c *Character) handleRollCharacterPhysical(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: c.GenerateAttributSelect(s, i),
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
}

func (c *Character) handleRollCharacterMental(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Rolling character")
	rolls, err := rollAttributes()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	log.Println(rolls)
	components := make([]discordgo.SelectMenuOption, len(rolls))
	for idx, roll := range rolls {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%v  %d", roll.Details, roll.Roll),
			Value: fmt.Sprintf("roll:%d:%d", roll.Index, roll.Roll),
		}
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Assign your physical attributes:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							Placeholder: "INT",
							CustomID:    fmt.Sprintf("%s:int", selectAttributeKey),
							Options:     components,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							Placeholder: "WIS",
							CustomID:    fmt.Sprintf("%s:wis", selectAttributeKey),
							Options:     components,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							Placeholder: "CHA",
							CustomID:    fmt.Sprintf("%s:char", selectAttributeKey),
							Options:     components,
						},
					},
				},
			},
		},
	}

	log.Println("response built")
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
	log.Println("response sent")
}

func (c *Character) handleRandomStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("handleRandomStart called")
	choices, err := c.startNewChoices(4)
	if err != nil {
		log.Println(err)
		// TODO: Handle error
		return
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

	rand.Seed(time.Now().UnixNano())
	for idx := 0; idx < 4; idx++ {
		choices[idx] = &charChoice{
			Race:  races[rand.Intn(len(races))],
			Class: classes[rand.Intn(len(classes))],
		}
	}

	return choices, nil
}

type rollResult struct {
	Index   int
	Roll    int
	Details []int
}

func rollAttributes() ([]*rollResult, error) {
	attributes := make([]*rollResult, 6)

	for idx := range attributes {
		roll, err := dice.Roll("4d6")
		if err != nil {
			return nil, err
		}
		attributes[idx] = &rollResult{
			Index:   idx,
			Roll:    roll.Total - roll.Lowest,
			Details: roll.Rolls,
		}
	}
	return attributes, nil
}
