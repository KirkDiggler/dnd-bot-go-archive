package components

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"log"
	"log/slog"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const ronnieRollBack = "ronnie-roll-back"

type RonnieD struct {
	messageID string
	manager   ronnied_actions.Interface
}

type RonnieDConfig struct {
	Manager ronnied_actions.Interface
}

func NewRonnieD(cfg *RonnieDConfig) (*RonnieD, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Manager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Manager")
	}

	return &RonnieD{
		manager: cfg.Manager,
	}, nil
}

func (c *RonnieD) RollBack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	oldInteraction := &discordgo.Interaction{AppID: i.AppID, Token: c.messageID}
	err := s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Print(err)
	}

	c.RonnieRoll(s, i)
}

func (c *RonnieD) RonnieRolls(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// if the roll count is greater than 5 return invalid input response
	if i.ApplicationCommandData().Options[0].Options[0].IntValue() > 5 {
		c.returnResponseMessage("Whoa there 🤠, slow down. 5 rolls max!", s, i)
		return
	}

	numberOfRolls := i.ApplicationCommandData().Options[0].Options[0].IntValue()
	rolls := make([]int, numberOfRolls)
	for idx := int64(0); idx < numberOfRolls; idx++ {
		rolls[idx] = rand.Intn(6) + 1
	}

	slog.Info("Rolls", "rolls", rolls)

	msgBuilder := strings.Builder{}
	var response *discordgo.InteractionResponse
	c.messageID = i.Token

	gameResult, err := c.manager.AddRolls(context.Background(), &ronnied_actions.AddRollsInput{
		GameID:    i.ChannelID,
		PlayerID:  i.Member.User.ID,
		RollCount: int(numberOfRolls),
	})
	if err != nil {
		log.Print(err)

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		}
		err = s.InteractionRespond(i.Interaction, response)
		if err != nil {
			log.Print(err)
		}
	}

	msgBuilder.WriteString(fmt.Sprintf("%s rolled %d times\n", i.Member.User.Username, numberOfRolls))

	slog.Info("Game Result", "gameResult", gameResult)
	if gameResult != nil && gameResult.Success {
		for _, result := range gameResult.Results {
			if result == nil {
				slog.Warn("Result is nil")
				continue
			}

			if result.PlayerID == "" {
				slog.Warn("Missing playerID", "result", result)

				msgBuilder.WriteString("MISSING DATA\n")
				// return error if this happens. fail fast
				continue
			}

			msgBuilder.WriteString(fmt.Sprintf("🎲: **%d** ", result.Roll))
			// TODO: create grabbag from user input generate this list (load from file, seeding process?)
			bag := []string{"🍺", "🍻", "🍷", "🥃", "🍸", "🍹", "🍾", "🥂", "🥤", "🧉", "🧊", "🥛", "🍼", "☕", "🫖", "🍵", "🧃", "🥤", "🧋", "🍶", "🍺", "🍻", "🍷", "🥃", "🍸", "🍹", "🍾", "🥂", "🥤", "🧉", "🧊", "🥛", "🍼", "☕", "🫖", "🍵", "🧃", "🥤", "🧋", "🍶", "🍺", "🍻", "🍷", "🥃", "🍸", "🍹", "🍾", "🥂", "🥤", "🧉", "🧊", "🥛", "🍼", "☕", "🫖", "🍵", "🧃", "🥤", "🧋", "🍶", "🍺", "🍻", "🍷", "🥃", "🍸", "🍹", "🍾", "🥂", "🥤", "🧉", "🧊", "🥛", "🍼", "☕", "🫖", "🍵", "🧃", "🥤", "🧋", "🍶", "🍺", "🍻", "🍷", "🥃", "🍸", "🍹", "🍾", "🥂", "🥤", "🧉", "🧊", "🥛", "🍼", "☕", "🫖", "🍵", "🧃", "🥤", "🧋", "🍶", "🍺"}
			grabbed := bag[rand.Intn(len(bag))] // this will be unique per row

			switch result.Roll {
			case 1:
				msgBuilder.WriteString(grabbed)
			case 6:
				if result.AssignedTo == "" {
					slog.Warn("Missing assignedTo", "result", result)

					// TODO: move to constant
					msgBuilder.WriteString("sir... sir I am missing data (check logs)")
					continue
				}

				user, userErr := s.User(result.AssignedTo)
				if userErr != nil {
					log.Print(userErr)
					return
				}

				msgBuilder.WriteString(fmt.Sprintf("%s %sノ( ゜-゜ノ)", grabbed, user.Username))
			default:
				// respond with trumpet emoji
				msgBuilder.WriteString("*sad trumpet*")
			}

			msgBuilder.WriteString("\n")

		}

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		}
	} else {
		roll := rand.Intn(6) + 1

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
									Emoji: &discordgo.ComponentEmoji{
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
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) returnResponseMessage(responseMessage string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) RonnieRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roll := rand.Intn(6) + 1

	msgBuilder := strings.Builder{}
	var response *discordgo.InteractionResponse
	c.messageID = i.Token

	gameResult, err := c.manager.AddRoll(context.Background(), &ronnied_actions.AddRollInput{
		GameID:   i.ChannelID,
		PlayerID: i.Member.User.ID,
		Roll:     roll,
	})
	if err != nil {
		log.Print(err)

		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		}
		err = s.InteractionRespond(i.Interaction, response)
		if err != nil {
			log.Print(err)
		}
	}

	if gameResult != nil && gameResult.Success {
		user, userErr := s.User(gameResult.AssignedTo)
		if userErr != nil {
			log.Print(userErr)
			return
		}

		msgBuilder.WriteString(fmt.Sprintf("%s rolled a %d and the drink was assigned to %s",
			i.Member.User.Username, roll,
			user.Username))
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msgBuilder.String(),
			},
		}
	} else {
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
									Emoji: &discordgo.ComponentEmoji{
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
									Emoji: &discordgo.ComponentEmoji{
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
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) RonnieRollBack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString(fmt.Sprintf("%s rolled it back!", i.Member.User.Username))
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: msgBuilder.String(),
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Print(err)
	}
}

func (c *RonnieD) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	grabBag := []string{
		"You know it Bro!",
		"Right back at you",
		"This guy get's it",
		"💯",
	}

	result := grabBag[rand.Intn(len(grabBag))]
	if m.Content == "thanks ronnie" ||
		m.Content == "Thank's Ronnie" ||
		m.Content == "Thank's ronnie" ||
		m.Content == "thank's Ronnie" ||
		m.Content == "thank's ronnie" ||
		m.Content == "thanks ronnie d" ||
		m.Content == "thank's ronnie d" {
		_, err := s.ChannelMessageSend(m.ChannelID, result)
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "tanks ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Get a load of this guy")
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "there it is" {
		_, err := s.ChannelMessageSend(m.ChannelID, "It's right there")
		if err != nil {
			log.Print(err)
		}
	} else if m.Content == "comon ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, "You got this")
		if err != nil {
			log.Print(err)
		}
	}

	if m.Content == "how about you ronnie" || m.Content == "how about you ronnie d" {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm good Bro, thanks for asking")
		if err != nil {
			log.Print(err)
		}
	}

	if m.Content == "tanks ronnie" || m.Content == "tanks Ronnie" {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Get a load of %s. I'm not a tank, I'm a healer", m.Author.Username))
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "ronnied":
			switch i.ApplicationCommandData().Options[0].Name {
			case "gamejoin":
				c.JoinGame(s, i)
			case "gettab":
				c.GetTab(s, i)
			case "tabs":
				c.ListTabs(s, i)
			case "roll":
				c.RonnieRoll(s, i)
			case "rolls":
				c.RonnieRolls(s, i)
			case "drink":
				c.PayDrink(s, i)
			case "advise":
				grabBag := []string{
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: that's a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: pass a drink", i.Member.User.Username),
					fmt.Sprintf("%s asked Ronnie D for advice, Ronnie D says: social!", i.Member.User.Username),
				}

				result := grabBag[rand.Intn(len(grabBag))]

				log.Println(result)

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
			c.RollBack(s, i)
		}
	}
}

func (c *RonnieD) GetTab(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "gettab" {
		var msg string

		result, err := c.manager.GetTab(context.Background(), &ronnied_actions.GetTabInput{
			GameID:   i.ChannelID,
			PlayerID: i.Member.User.ID,
		})
		if err != nil {
			log.Print(err)
			msg = err.Error()
		} else {
			msg = fmt.Sprintf("Your tab is %d", result.Count)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) ListTabs(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	msg := strings.Builder{}

	if data.Options[0].Name == "tabs" {
		result, err := c.manager.ListTabs(context.Background(), &ronnied_actions.ListTabsInput{
			GameID: i.ChannelID,
		})
		if err != nil {
			log.Print(err)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			if err != nil {
				log.Print(err)
			}
			return
		}

		for _, tab := range result.Tabs {
			user, userErr := s.User(tab.PlayerID)
			if userErr != nil {
				log.Print(userErr)
			}

			msg.WriteString(fmt.Sprintf("Player: %s: %d\n", user.Username, tab.Count))
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg.String(),
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

// AddResult adds a result to a game
// TODO: move to the roll command.  All rolls should be sent and based on the success response we will send a message
func (c *RonnieD) AddResult(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "addresult" {
		gameID := data.Options[0].Options[0].StringValue()
		roll := data.Options[0].Options[1].IntValue()
		result, err := c.manager.AddRoll(context.Background(), &ronnied_actions.AddRollInput{
			GameID:   gameID,
			PlayerID: i.Member.User.ID,
			Roll:     int(roll),
		})
		if err != nil {
			log.Print(err)
			return
		}

		if result.Success {
			msg := fmt.Sprintf("Drink assigned to %s", result.AssignedTo)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
			if err != nil {
				log.Print(err)
			}
		}

	}
}

func (c *RonnieD) PayDrink(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gameID := i.ChannelID
	// Get the channel name

	builder := strings.Builder{}
	builder.WriteString("Prepare to drink \n\n...\n\n")

	_, err := c.manager.PayDrink(context.Background(), &ronnied_actions.PayDrinkInput{
		GameID:   gameID,
		PlayerID: i.Member.User.ID,
	})
	if err != nil {
		builder.WriteString(err.Error())
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: builder.String(),
		},
	})
	if err != nil {
		log.Print(err)
	}

}

func (c *RonnieD) JoinGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "gamejoin" {
		gameID := i.ChannelID
		// Get the channel name
		channel, err := s.Channel(i.ChannelID)
		if err != nil {
			log.Print(err)
			return
		}
		msg := fmt.Sprintf("You joined the game")

		_, err = c.manager.JoinGame(context.Background(), &ronnied_actions.JoinGameInput{
			GameID:   gameID,
			GameName: channel.Name,
			PlayerID: i.Member.User.ID,
		})
		if err != nil {
			msg = err.Error()
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
		}
	}
}

func (c *RonnieD) CreateGame(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Options[0].Name == "creategame" {
		gameName := data.Options[0].Options[0].StringValue()

		msg := fmt.Sprintf("Game %s created, ID: %d", gameName, rand.Intn(1000))
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Print(err)
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
				Name:        "rolls",
				Description: "Rolls for RonnieD",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "number",
						Description: "How many rolls to make",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Required:    true,
					},
				},
			}, {
				Name:        "advise",
				Description: "what should I do RonnieD?",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "gamejoin",
				Description: "Join a game",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "gettab",
				Description: "Get your tab",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "tabs",
				Description: "List all tabs for this channel",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "drink",
				Description: "drink your tab",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}
