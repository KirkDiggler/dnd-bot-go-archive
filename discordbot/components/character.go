package components

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	entities2 "github.com/KirkDiggler/dnd-bot-go/internal/entities"
	character2 "github.com/KirkDiggler/dnd-bot-go/internal/repositories/character"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/bwmarrin/discordgo"
)

const selectCaracterAction = "select-character"

type Character struct {
	client   dnd5e.Interface
	charRepo character2.Repository
}

type CharacterConfig struct {
	Client        dnd5e.Interface
	CharacterRepo character2.Repository
}

type charChoice struct {
	Name  string
	Race  *entities2.Race
	Class *entities2.Class
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
		client:   cfg.Client,
		charRepo: cfg.CharacterRepo,
	}, nil
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
		}
	}
}
func (c *Character) handleLoadCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charRepo.GetCharacter(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Loaded character %s the %s %s", char.Name, char.RaceKey, char.ClassKey),
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

	char, err := c.charRepo.CreateCharacter(context.Background(), &character2.Data{
		OwnerID:  i.Member.User.ID,
		Name:     i.Member.User.Username,
		RaceKey:  race,
		ClassKey: class,
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
		},
	})
	if err != nil {
		log.Println(err)
		return
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
