package components

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type RonnieD struct {
}

func NewRonnieD() (*RonnieD, error) {
	return &RonnieD{}, nil
}

func (c *RonnieD) RonnieRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(6) + 1
	msgBuilder := strings.Builder{}
	if roll == 6 {
		msgBuilder.WriteString("Crit! Pass a drink")
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.SuccessButton,
								CustomID: "roll",
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Println(err)
		}
	} else if roll == 1 {
		msgBuilder.WriteString("rolled a 1, that's a drink!")
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Roll it back",
								Style:    discordgo.PrimaryButton,
								CustomID: "ronnied-roll-back",
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
		}
	} else {
		msgBuilder.WriteString("You rolled a ")
		msgBuilder.WriteString(strconv.Itoa(roll))
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		})
		if err != nil {
			log.Println(err)
		}
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
					"Ronnie D says: That's a drink",
					"Ronnie D says: Pass a drink",
					"Ronnie D says: Social!",
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
		case "ronnied-roll-back":
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
