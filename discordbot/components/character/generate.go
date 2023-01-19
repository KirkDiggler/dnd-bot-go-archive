package character

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
	Name  string
	Race  *entities.Race
	Class *entities.Class
}

type Component struct {
	client        dnd5e.Interface
	session       *discordgo.Session
	appID         string
	guildID       string
	choices       []*charChoice
	choiceOptions *choiceOptions
	poll          *poll.Poll
}

type Config struct {
	Client  dnd5e.Interface
	Session *discordgo.Session
	AppID   string
	GuildID string
}

type charChoice struct {
	Name  string
	Race  *entities.Race
	Class *entities.Class
}

type choiceOptions struct {
	Names   []string
	Races   []*entities.Race
	Classes []*entities.Class
}

func newChoiceOptions() *choiceOptions {
	return &choiceOptions{
		Names:   make([]string, 0),
		Races:   make([]*entities.Race, 0),
		Classes: make([]*entities.Class, 0),
	}
}

func New(cfg *Config) (*Component, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.Session == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Session")
	}

	if cfg.AppID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.AppID")
	}

	if cfg.GuildID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.GuildID")
	}

	return &Component{
		client:        cfg.Client,
		session:       cfg.Session,
		appID:         cfg.AppID,
		guildID:       cfg.GuildID,
		choiceOptions: newChoiceOptions(),
		poll:          poll.New(),
	}, nil
}

func (c *Component) startNewChoices(number int) {
	fmt.Println("Starting new choices")
	choices := make([]*charChoice, number)
	rand.Seed(time.Now().UnixNano())
	for idx := 0; idx < 4; idx++ {
		choices[idx] = &charChoice{
			Name:  c.choiceOptions.Names[rand.Intn(len(c.choiceOptions.Names))],
			Race:  c.choiceOptions.Races[rand.Intn(len(c.choiceOptions.Races))],
			Class: c.choiceOptions.Classes[rand.Intn(len(c.choiceOptions.Classes))],
		}
	}

	c.choices = choices
	c.poll = poll.New()
}

func (c *Component) Load(s *discordgo.Session) error {
	races, err := c.client.ListRaces()
	if err != nil {
		return err
	}

	classes, err := c.client.ListClasses()
	if err != nil {
		return err
	}

	allMembers, err := s.GuildMembers(c.guildID, "", 1000)
	if err != nil {
		log.Println(err)
	}

	possibleMembers := make([]string, 0)
	log.Printf("%d members loaded", len(allMembers))
	for _, member := range allMembers {
		if !member.User.Bot {
			possibleMembers = append(possibleMembers, member.User.Username)
		}
	}

	c.choiceOptions = &choiceOptions{
		Names:   possibleMembers,
		Races:   races,
		Classes: classes,
	}

	if len(possibleMembers) < 4 {
		log.Println("Not enough members to generate a party")
	}

	c.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "generate":
				c.GenerateStartHandler()(s, i)
			}
		case discordgo.InteractionMessageComponent:
			switch i.MessageComponentData().CustomID {
			case "vote-choice":
				c.VoteChoiceHandler()(s, i)
			}
		}
	})

	_, err = c.session.ApplicationCommandCreate(c.appID, c.guildID, &discordgo.ApplicationCommand{
		Name:        "generate",
		Description: "Generate a random list of characters to choose from",
	})

	return err
}

func (c *Component) VoteChoiceHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
}

func (c *Component) GenerateStartHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	c.startNewChoices(4)

	components := make([]discordgo.SelectMenuOption, len(c.choices))
	for idx, char := range c.choices {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%s the %s %s", char.Name, char.Race.Name, char.Class.Name),
			Value: fmt.Sprintf("choice-%d", idx),
		}
	}

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		println("generate start handler")
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
}
