package components

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	selectCaracterAction = "select-character"
	rollCharacterAction  = "roll-character"
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
				Description: "Create a new character from a random list",
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
		switch i.MessageComponentData().CustomID {
		case selectCaracterAction:
			c.handleCharSelect(s, i)
		case rollCharacterAction:
			c.handleRollCharacter(s, i)
		}
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

	char, err := c.charManager.Create(context.Background(), &entities.Character{
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
							CustomID: rollCharacterAction,
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

func (c *Character) handleRollCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rolls, err := rollAttributes()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	components := make([]discordgo.SelectMenuOption, len(rolls))
	for idx, roll := range rolls {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%v  %d", roll.Details, roll.Roll),
			Value: fmt.Sprintf("roll:%d:%d", roll.Index, roll.Roll),
		}
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Assign your rolls to your stats",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							Placeholder: "DEX",
							CustomID:    fmt.Sprintf("%s:DEX", selectCaracterAction),
							Options:     components,
						},
						&discordgo.SelectMenu{
							Placeholder: "STR",
							CustomID:    fmt.Sprintf("%s:STR", selectCaracterAction),
							Options:     components,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							Placeholder: "WIS",
							CustomID:    fmt.Sprintf("%s:WIS", selectCaracterAction),
							Options:     components,
						},
						&discordgo.SelectMenu{
							Placeholder: "INT",
							CustomID:    fmt.Sprintf("%s:INT", selectCaracterAction),
							Options:     components,
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.SelectMenu{
									Placeholder: "CON",
									CustomID:    fmt.Sprintf("%s:CON", selectCaracterAction),
									Options:     components,
								},
								&discordgo.SelectMenu{
									Placeholder: "CHAR",
									CustomID:    fmt.Sprintf("%s:CHAR", selectCaracterAction),
									Options:     components,
								},
							},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "Save",
							CustomID: fmt.Sprintf("%s:save", selectCaracterAction),
						},
					},
				},
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
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
			Details: roll.Details,
		}
	}
	return attributes, nil
}
