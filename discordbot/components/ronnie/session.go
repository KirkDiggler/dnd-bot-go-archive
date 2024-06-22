package ronnie

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strings"
)

const socialGameID = "ronnie-rollem"

//func (c *RonnieD) handleSessionNew(s *discordgo.Session, i *discordgo.InteractionCreate, sessionRollID string) {
//	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
//		SessionRollID: sessionRollID,
//	})
//	if err != nil {
//		fmt.Println("Failed to get session roll:", err)
//		return
//	}
//
//	//session, err := c.manager.GetSession(context.Background(), &ronnied_actions.GetSessionInput{
//	//	SessionID: sessionRollResult.SessionRoll.SessionID,
//	//})
//	//if err != nil {
//	//	fmt.Println("Failed to get session:", err)
//	//	return
//	//}
//	//
//	//err = s.ChannelMessageDelete(i.ChannelID, session.Session.MessageID)
//	//if err != nil {
//	//	log.Println("Failed to delete message:", err)
//	//}
//
//	c.SessionCreate(s, i)
//}

//func (c *RonnieD) SessionNew(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	// Cleanup the old messages and reset the exiusting message for the new session
//	// Get the session ID from the CustomID
//	data := i.MessageComponentData()
//	sessionRollID := data.CustomID[len("new_session:"):]
//
//	log.Println("sessionRollID", sessionRollID)
//	c.handleSessionNew(s, i, sessionRollID)
//}

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

func (c *RonnieD) sessionJoin(s *discordgo.Session, i *discordgo.InteractionCreate, sessionRollID string) {
	// TODO: use session id and find
	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	user, err := s.User(i.Member.User.ID)
	if err != nil {
		log.Println("Failed to get user:", err)
		return
	}

	// Join the session
	sessionResult, err := c.manager.JoinSession(context.Background(), &ronnied_actions.JoinSessionInput{
		SessionID:     sessionRollResult.SessionRoll.SessionID,
		SessionRollID: sessionRollResult.SessionRoll.ID,
		PlayerID:      i.Member.User.ID,
		PlayerName:    user.Username,
	})
	if err != nil {
		fmt.Println("Failed to join session:", err)
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to join the session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// Respond with a message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Joined game with ID %s.", getLastPartID(sessionResult.Session.ID)),
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Rollem",
						CustomID: "rollem:" + sessionRollID,
						Style:    discordgo.PrimaryButton,
					},
				}},
			},
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}

	c.updateGameMessage(s, i, sessionRollID)
}
func (c *RonnieD) SessionJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the session ID from the CustomID
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("join_session:"):]
	log.Println("sessionRollID", sessionRollID)
	c.sessionJoin(s, i, sessionRollID)
}

func (c *RonnieD) SessionJoinReserved(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the session ID from the CustomID
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("join_reserved_session:"):]
	log.Println("sessionRollID", sessionRollID)
	c.sessionJoin(s, i, sessionRollID)
}

// SessionContinue marks your entry as complete and updates the main game message
func (c *RonnieD) SessionContinue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("session_continue:"):]
	log.Println("sessionRollID", sessionRollID)

	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	// Check if the player is part of the session
	if player := sessionRollResult.SessionRoll.HasPlayer(i.Member.User.ID); player == nil {
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not part of this session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// Find their entry and mark it completed
	entry := sessionRollResult.SessionRoll.HasPlayerEntry(i.Member.User.ID)
	if entry == nil {
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You have not rolled yet.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	entry.Completed = true

	_, err = c.manager.UpdateSessionRoll(context.Background(), &ronnied_actions.UpdateSessionRollInput{
		SessionRoll: sessionRollResult.SessionRoll,
	})
	if err != nil {
		fmt.Println("Failed to update session roll:", err)
	}

	// update the message with the new status
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "You have completed your turn. " + entry.String(),
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}

	c.updateGameMessage(s, i, sessionRollID)
}

func (c *RonnieD) SessionRoll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("rollem:"):]
	log.Println("sessionRollID", sessionRollID)

	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	// Check if the player is part of the session
	if player := sessionRollResult.SessionRoll.HasPlayer(i.Member.User.ID); player == nil {
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not part of this session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	var content string

	rollResult, err := c.manager.AddSessionRoll(context.Background(), &ronnied_actions.AddSessionRollInput{
		SessionRollID: sessionRollResult.SessionRoll.ID,
		PlayerID:      i.Member.User.ID,
	})
	if err != nil {
		fmt.Println("c.manager.AddSessionRoll returned err:", err)
		content = err.Error()
	} else {
		content = fmt.Sprintf("You rolled a %d.", rollResult.SessionEntry.Roll)
	}

	if rollResult != nil {
		if rollResult.SessionEntry.Roll == 1 {
			content = content + " That's a drink! 🍻"
			_, err = c.manager.AssignDrink(context.Background(), &ronnied_actions.AssignDrinkInput{
				GameID:        socialGameID,
				SessionRollID: sessionRollID,
				PlayerID:      i.Member.User.ID,
				AssignedTo:    i.Member.User.ID,
			})
			if err != nil {
				fmt.Println("Failed to assign drink:", err)
				content = content + " Failed to assign drink."
			}
		}

		if rollResult.SessionEntry.Roll == 6 {
			// Add a dropdown with the session's players to assign the drink to
			// we respond here so the dont see the continue button yet
			var options []discordgo.SelectMenuOption
			for _, player := range sessionRollResult.SessionRoll.Players {
				user, userErr := s.User(player.ID)
				if userErr != nil {
					log.Println("Failed to get user:", userErr)
					continue
				}

				options = append(options, discordgo.SelectMenuOption{
					Label: user.Username,
					Value: player.ID,
				})
			}

			// Respond with the roll and dropdown
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: content,
					Components: []discordgo.MessageComponent{
						&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
							&discordgo.SelectMenu{
								Placeholder: "Pass the drink to",
								CustomID:    "assign_drink:" + sessionRollID,
								Options:     options,
							},
						}},
					},
				},
			})
			if err != nil {
				log.Println("Failed to respond to interaction:", err)
				return
			}
			return
		}
	}

	// Respond with the roll and continue buttone
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}

	c.updateGameMessage(s, i, sessionRollID)
}

func (c *RonnieD) SessionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	result, err := c.manager.CreateSession(context.Background(), &ronnied_actions.CreateSessionInput{
		GameID: socialGameID,
	})
	if err != nil {
		fmt.Println("Failed to create session:", err)
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create a session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	sessionRoll, err := c.manager.CreateSessionRoll(context.Background(), &ronnied_actions.CreateSessionRollInput{
		SessionID:    result.Session.ID,
		Type:         ronnied.RollTypeStart,
		Participants: []string{i.Member.User.ID},
	})
	if err != nil {
		fmt.Println("Failed to create session roll:", err)
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create a session roll.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// joing the user to the new session
	user, err := s.User(i.Member.User.ID)
	if err != nil {
		log.Println("Failed to get user:", err)
		return
	}
	// Join the session
	_, err = c.manager.JoinSession(context.Background(), &ronnied_actions.JoinSessionInput{
		SessionID:     sessionRoll.SessionRoll.SessionID,
		SessionRollID: sessionRoll.SessionRoll.ID,
		PlayerID:      i.Member.User.ID,
		PlayerName:    user.Username,
	})
	if err != nil {
		fmt.Println("Failed to join session:", err)
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to join the session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// Create an embed for the session
	msgID, err := c.sendGameStartMessage(s, i, sessionRoll.SessionRoll.ID)
	if err != nil {
		fmt.Println("Failed to send game start message:", err)
		return
	}

	result.Session.MessageID = msgID
	// Save the message ID to the session
	_, err = c.manager.UpdateSession(context.Background(), &ronnied_actions.UpdateSessionInput{
		Session: result.Session,
	})
	if err != nil {
		fmt.Println("Failed to update session:", err)
	}
}

func (c *RonnieD) SessionAssignDrink(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("assign_drink:"):]
	log.Println("sessionRollID", sessionRollID)

	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	// Check if the player is part of the session
	if player := sessionRollResult.SessionRoll.HasPlayer(i.Member.User.ID); player == nil {
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not part of this session.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// Get the assigned player ID from the dropdown
	assignedPlayerID := data.Values[0]
	log.Println("assignedPlayerID", assignedPlayerID)

	// Assign the drink
	_, err = c.manager.AssignDrink(context.Background(), &ronnied_actions.AssignDrinkInput{
		SessionRollID: sessionRollID,
		PlayerID:      i.Member.User.ID,
		AssignedTo:    assignedPlayerID,
	})
	if err != nil {
		fmt.Println("Failed to assign drink:", err)
		respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to assign the drink.",
			},
		})
		if respondErr != nil {
			log.Println("Failed to respond to interaction:", respondErr)
		}

		return
	}

	// Respond with a message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("You assigned the drink to <@%s>.", assignedPlayerID),
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}

	c.updateGameMessage(s, i, sessionRollID)
}

type results struct {
	player string
	drinks int
}

type sessionRollRenderInput struct {
	SessionRollID string
	RollType      ronnied.RollType
	Session       *discordgo.Session
}

func (c *RonnieD) renderSessionRoll(input *sessionRollRenderInput) (*discordgo.MessageEmbed, error) {
	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: input.SessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return nil, err
	}

	// Get the session details
	result, err := c.manager.GetSession(context.Background(), &ronnied_actions.GetSessionInput{
		SessionID: sessionRollResult.SessionRoll.SessionID,
	})
	if err != nil {
		fmt.Println("Failed to get session:", err)
		return nil, err
	}
	sessionID := getLastPartID(result.Session.ID)

	// set the main title based on the type
	var title string
	switch input.RollType {
	case ronnied.RollTypeRollOff:
		title = fmt.Sprintf("Game %s roll off", sessionID)
	default:
		title = fmt.Sprintf("Game %s accepting players", sessionID)
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0x00ff00, // Green color
	}

	playerResults := make(map[string]*results)

	for _, participant := range sessionRollResult.SessionRoll.Players {
		playerResults[participant.ID] = &results{
			drinks: 0,
			player: participant.Name,
		}
	}

	for _, participant := range sessionRollResult.SessionRoll.Players {
		// add green checkmark icon if the entry is completed
		var content string

		addDefault := true
		if entry := sessionRollResult.SessionRoll.HasPlayerEntry(participant.ID); entry != nil {
			embed.Title = fmt.Sprintf("Rolling %s", sessionID)
			//default to a timeer icon
			content = "⏳ "
			addDefault = false
			if entry.Completed {
				content = "✅ "
			}
			if sessionRollResult.SessionRoll.IsComplete() {
				if sessionRollResult.SessionRoll.IsLoser(entry) {
					content = "🍻 "
					playerResults[participant.ID].drinks += 1
				} else {
					content = "🎉 "
				}
			}

			if entry.AssignedTo != "" {
				if _, ok := playerResults[entry.AssignedTo]; !ok {
					log.Println("Assigned user not found:", entry.AssignedTo, " adding to results")
					playerResults[entry.AssignedTo] = &results{
						player: entry.AssignedTo,
						drinks: 0,
					}
				}

				playerResults[entry.AssignedTo].drinks += 1
				assignedUser, assignedUserErr := input.Session.User(entry.AssignedTo)
				if assignedUserErr != nil {
					log.Println("Failed to get user:", assignedUserErr)
					continue
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   participant.Name,
					Value:  content + fmt.Sprintf("🎲 %d, 🍻 %s", entry.Roll, assignedUser.Username),
					Inline: true,
				})
				addDefault = false
				continue
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   participant.Name,
				Value:  content + fmt.Sprintf("🎲 %d", entry.Roll),
				Inline: true,
			})
		}

		if sessionRollResult.SessionRoll.IsComplete() {
			embed.Title = "Rolled"
		}

		if addDefault {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   participant.Name,
				Value:  "waiting for roll",
				Inline: true,
			})
		}
	}

	return embed, nil
}
func (c *RonnieD) updateGameMessage(s *discordgo.Session, i *discordgo.InteractionCreate, sessionRollID string) {
	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	// Get the session details
	result, err := c.manager.GetSession(context.Background(), &ronnied_actions.GetSessionInput{
		SessionID: sessionRollResult.SessionRoll.SessionID,
	})
	if err != nil {
		fmt.Println("Failed to get session:", err)
		return
	}

	embed, err := c.renderSessionRoll(&sessionRollRenderInput{
		SessionRollID: sessionRollID,
		Session:       s,
	})
	if err != nil {
		fmt.Println("Failed to render session roll:", err)
		return
	}

	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "New Game",
			CustomID: "new_session:" + sessionRollResult.SessionRoll.ID, // Include session ID in CustomID
			Style:    discordgo.PrimaryButton,
		},
	}

	playerResults := make(map[string]*results)

	for _, participant := range sessionRollResult.SessionRoll.Players {
		playerResults[participant.ID] = &results{
			drinks: 0,
			player: participant.Name,
		}
	}

	embeds := []*discordgo.MessageEmbed{embed}

	if sessionRollResult.SessionRoll.IsComplete() {
		embed.Color = 0xff0000 // Red color
		//add results with yellow beer color
		resultEmbed := &discordgo.MessageEmbed{
			Title: "Dranks! 🍻",
			Color: 0x00ff00, // Green color
		}

		for _, playerResult := range playerResults {
			resultEmbed.Fields = append(resultEmbed.Fields, &discordgo.MessageEmbedField{
				Name:   playerResult.player,
				Value:  fmt.Sprintf("🍻 %d", playerResult.drinks),
				Inline: true,
			})
		}

		embeds = append(embeds, resultEmbed)

		winners := sessionRollResult.SessionRoll.GetWinners()

		if len(winners) > 1 {
			playerIDs := make([]string, len(winners))
			for idx, winner := range winners {
				playerIDs[idx] = winner.PlayerID
			}

			// create new session roll with the winners and render the session roll
			newSessionRoll, rollOffErr := c.manager.CreateSessionRoll(context.Background(), &ronnied_actions.CreateSessionRollInput{
				SessionID:    sessionRollResult.SessionRoll.SessionID,
				Type:         ronnied.RollTypeStart,
				Participants: playerIDs,
			})
			if rollOffErr != nil {
				fmt.Println("Failed to create session roll:", rollOffErr)
				return
			}

			rollOffEmbed, rollOffErr := c.renderSessionRoll(&sessionRollRenderInput{
				SessionRollID: newSessionRoll.SessionRoll.ID,
				RollType:      ronnied.RollTypeRollOff,
				Session:       s,
			})
			if rollOffErr != nil {
				fmt.Println("Failed to render session roll:", rollOffErr)
				return
			}

			embeds = append(embeds, rollOffEmbed)
			buttons = append(buttons,
				discordgo.Button{
					Label:    "Join Game",
					CustomID: "join_reserved_session:" + newSessionRoll.SessionRoll.ID, // Include session ID in CustomID
					Style:    discordgo.SuccessButton,
				},
			)
		}
	}

	if !sessionRollResult.SessionRoll.IsComplete() {
		buttons = append(buttons,
			discordgo.Button{
				Label:    "Join Game",
				CustomID: "join_session:" + sessionRollResult.SessionRoll.ID, // Include session ID in CustomID
				Style:    discordgo.SuccessButton,
			},
		)
	}

	buttonRow := []discordgo.MessageComponent{
		&discordgo.ActionsRow{Components: buttons},
	}

	messageTitle := "Roll em!"
	edit := &discordgo.MessageEdit{
		ID:         result.Session.MessageID,
		Content:    &messageTitle,
		Channel:    i.ChannelID,
		Embeds:     &embeds,
		Components: &buttonRow,
	}

	_, err = s.ChannelMessageEditComplex(edit)
	if err != nil {
		log.Println("Failed to edit message:", err)

		return
	}
}

func (c *RonnieD) sendGameStartMessage(s *discordgo.Session, i *discordgo.InteractionCreate, sessionRollID string) (string, error) {
	rollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return "", err
	}

	embed := &discordgo.MessageEmbed{
		Title:       "New Rollem Game started",
		Description: "Waiting for all players to join",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  "Waiting for players...",
				Inline: true,
			},
		},
		Color: 0x00ff00, // Green color
	}

	for _, participant := range rollResult.SessionRoll.Players {
		user, userErr := s.User(participant.ID)
		if userErr != nil {
			log.Println("Failed to get user:", userErr)
			continue
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Player",
			Value:  user.Username,
			Inline: true,
		})
	}

	// Create a join button
	joinButton := discordgo.Button{
		Label:    "Prepare to roll",
		CustomID: "join_session:" + rollResult.SessionRoll.ID, // Include session ID in CustomID
		Style:    discordgo.SuccessButton,
	}

	newButton := discordgo.Button{
		Label:    "New Game",
		CustomID: "new_session:" + rollResult.SessionRoll.ID, // Include session ID in CustomID
		Style:    discordgo.PrimaryButton,
	}
	// Respond with the embed and button
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{newButton, joinButton}},
			},
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Println("Failed to get interaction response:", err)
		return "", err
	}

	return msg.ID, nil
}

func getLastPartID(id string) string {
	parts := strings.Split(id, "-")

	return parts[len(parts)-1]
}
