package components

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/entities"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/repositories/party"
	"github.com/bwmarrin/discordgo"
)

type Party struct {
	// contains filtered or unexported fields
	appID             string
	guildID           string
	session           *discordgo.Session
	partyRepo         party.Interface
	componentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	commandHandlers   map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type PartyConfig struct {
	AppID     string
	GuildID   string
	Session   *discordgo.Session
	PartyRepo party.Interface
}

func NewParty(cfg *PartyConfig) (*Party, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.PartyRepo == nil {
		return nil, dnderr.NewMissingParameterError("cfg.PartyRepo")
	}

	if cfg.Session == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Session")
	}

	return &Party{
		session:           cfg.Session,
		partyRepo:         cfg.PartyRepo,
		componentHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		commandHandlers:   getCommandHandlers(),
	}, nil
}

func getCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"party": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.ApplicationCommandData().Options[0].Name {
			case "create":
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseModal,
					Data: &discordgo.InteractionResponseData{
						CustomID: "party-create-" + i.Interaction.Member.User.ID,
						Title:    "Create a new party",
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.TextInput{
										CustomID:    "parrty-name",
										Label:       "What is your name for the party?",
										Style:       discordgo.TextInputShort,
										Placeholder: "I thought you said Mongos",
										Required:    true,
										MaxLength:   50,
										MinLength:   3,
									},
								},
							},
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.TextInput{
										CustomID: "parrty-size",
										Label:    "What is the maximum party size?",
										Style:    discordgo.TextInputShort,
										Value:    "2",
										Required: true,
									},
								},
							},
						},
					},
				})
				if err != nil {
					fmt.Println(err)
					//panic(err)
				}
			}
		},
	}
}

func (c *Party) GetApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "party",
		Description: "Manage your party",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a new party",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c *Party) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if handler, ok := c.commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}

	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()
		msgBuilder := strings.Builder{}

		partySize := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		partyName := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		intSize, err := strconv.Atoi(partySize)
		if err != nil {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Invalid party size (must be a number)",
				},
			})
			if err != nil {
				panic(err)
			}

			return
		}

		partyData, err := c.partyRepo.CreateParty(context.Background(), &entities.Party{
			Name:      partyName,
			PartySize: intSize,
		})

		if err != nil {
			fmt.Println(err)
			msgBuilder.WriteString("Error creating party ")
		} else {
			msgBuilder.WriteString("Created party\n")
			msgBuilder.WriteString(fmt.Sprintf("Name: %s\n", partyData.Name))
			msgBuilder.WriteString(fmt.Sprintf("Size: %d\n", partyData.PartySize))
			msgBuilder.WriteString(fmt.Sprintf("Token: %s\n", partyData.Token))
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		})
		if err != nil {
			panic(err)
		}
	}
}
