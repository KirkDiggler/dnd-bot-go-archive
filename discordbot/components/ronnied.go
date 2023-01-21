package components

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const ronnieRollBack = "ronnie-roll-back"

type RonnieD struct {
}

func NewRonnieD() (*RonnieD, error) {
	return &RonnieD{}, nil
}

func (c *RonnieD) RonnieRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(6) + 1
	msgBuilder := strings.Builder{}
	var response *discordgo.InteractionResponse
	if roll == 6 {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a Crit! Pass a drink", i.Member.User.Username))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.SuccessButton,
								CustomID: ronnieRollBack,
								Emoji: discordgo.ComponentEmoji{
									Name: "🍺",
								},
							},
						},
					},
				},
			},
		}
	} else if roll == 1 {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a 1, that's a drink!", i.Member.User.Username))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.DangerButton,
								CustomID: ronnieRollBack,
								Emoji: discordgo.ComponentEmoji{
									Name: "🍺",
								},
							},
						},
					},
				},
			},
		}
	} else {
		msgBuilder.WriteString(fmt.Sprintf("%s rolled a %d, try again", i.Member.User.Username, roll))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		}
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ronnied":
			switch i.ApplicationCommandData().Options[0].Name {
			case "roll":
				c.RonnieRoll(s, i)
			case "advise":
				grabBag := []string{
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: that's a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: pass a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: social!", i.Member.User.Username),
				}

				result := grabBag[rand.Intn(len(grabBag)-1)]

				log.Printf("%s used ronnied and got %s", i.Member.User.Username, result)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: result,
					},
				})
				if err != nil {
					log.Print(err)
				}
			}
		}
	case discordgo.InteractionMessageComponent:
		switch i.MessageComponentData().CustomID {
		case ronnieRollBack:
			c.RonnieRoll(s, i)
		}
	}
}

func (c *RonnieD) GetApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ronnied",
		Description: "Leverage RonnieD's wisdom",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "roll",
				Description: "roll em!",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "advise",
				Description: "what should I do RonnieD?",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}
