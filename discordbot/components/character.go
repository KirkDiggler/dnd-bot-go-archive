package components

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/discordbot/components/poll"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/entities"
	"github.com/bwmarrin/discordgo"
)

type Character struct {
	client dnd5e.Interface

	choices       []*charChoice
	choiceOptions *choiceOptions
	poll          *poll.Poll
}

type CharacterConfig struct {
	Client dnd5e.Interface
}

type charChoice struct {
	Name  string
	Race  *entities.Race
	Class *entities.Class
}

type choiceOptions struct {
	Races   []*entities.Race
	Classes []*entities.Class
}

func newChoiceOptions() *choiceOptions {
	return &choiceOptions{
		Races:   make([]*entities.Race, 0),
		Classes: make([]*entities.Class, 0),
	}
}

func NewCharacter(cfg *CharacterConfig) (*Character, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	return &Character{
		client:        cfg.Client,
		choiceOptions: newChoiceOptions(),
		poll:          poll.New(),
	}, nil
}

func (c *Character) startNewChoices(number int) error {
	races, err := c.client.ListRaces()
	if err != nil {
		return err
	}

	classes, err := c.client.ListClasses()
	if err != nil {
		return err
	}
	fmt.Println("Starting new choices")
	choices := make([]*charChoice, number)

	rand.Seed(time.Now().UnixNano())
	for idx := 0; idx < 4; idx++ {
		choices[idx] = &charChoice{
			Race:  races[rand.Intn(len(races))],
			Class: classes[rand.Intn(len(classes))],
		}
	}

	c.choices = choices
	c.poll = poll.New()

	return nil
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
			}
		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case "vote-choice":
			c.handleVote(s, i)
		}
	}
}

func (c *Character) handleVote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	idxstring := strings.Replace(i.MessageComponentData().Values[0], "choice-", "", 1)
	idx, err := strconv.Atoi(idxstring)
	if err != nil {
		log.Println(err)
		return
	}

	if len(c.choices) < idx {
		log.Println("Index out of range")
		return
	}

	c.poll.Vote(i.Member.User.ID, idx)
	pollResults := c.poll.GetVotes()
	msgBuilder := strings.Builder{}
	println(fmt.Sprintf("Total votes: %v", pollResults))

	for idx, choice := range c.choices {
		vote := 0
		if v, ok := pollResults[idx]; ok {
			vote = v
		}

		msgBuilder.WriteString(fmt.Sprintf("%s the %s %s %d votes\n", choice.Name, choice.Race.Name, choice.Class.Name, vote))
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msgBuilder.String(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

}

func (c *Character) handleRandomStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := c.startNewChoices(4)
	if err != nil {
		log.Println(err)
		// TODO: Handle error
		return
	}

	components := make([]discordgo.SelectMenuOption, len(c.choices))
	for idx, char := range c.choices {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%s the %s %s", i.Member.User.Username, char.Race.Name, char.Class.Name),
			Value: fmt.Sprintf("choice-%d", idx),
		}
	}

	println("generate start handler")
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Place your vote for the next character:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    "vote-choice",
							Placeholder: "This is the tale of...👇",
							Options:     components,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}

}
