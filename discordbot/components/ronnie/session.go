package ronnie

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities/ronnied"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"github.com/bwmarrin/discordgo"
	"log"
)

const socialGameID = "ronnie-rollem"

func (c *RonnieD) SessionNew(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Cleanup the old messages and reset the exiusting message for the new session
	// Get the session ID from the CustomID
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("new_session:"):]

	log.Println("sessionRollID", sessionRollID)

	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	session, err := c.manager.GetSession(context.Background(), &ronnied_actions.GetSessionInput{
		SessionID: sessionRollResult.SessionRoll.SessionID,
	})
	if err != nil {
		fmt.Println("Failed to get session:", err)
		return
	}

	err = s.ChannelMessageDelete(i.ChannelID, session.Session.MessageID)
	if err != nil {
		log.Println("Failed to delete message:", err)
	}

	// get the players and delete the messages by MsgID
	for _, player := range sessionRollResult.SessionRoll.Players {
		err = s.ChannelMessageDelete(i.ChannelID, player.MsgID)
		if err != nil {
			log.Println("Failed to delete message:", err)
		}
	}

	c.SessionCreate(s, i)
}
func (c *RonnieD) SessionJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the session ID from the CustomID
	data := i.MessageComponentData()
	sessionRollID := data.CustomID[len("join_session:"):]
	log.Println("sessionRollID", sessionRollID)

	sessionRollResult, err := c.manager.GetSessionRoll(context.Background(), &ronnied_actions.GetSessionRollInput{
		SessionRollID: sessionRollID,
	})
	if err != nil {
		fmt.Println("Failed to get session roll:", err)
		return
	}

	// Join the session
	joinResult, err := c.manager.JoinSession(context.Background(), &ronnied_actions.JoinSessionInput{
		SessionID:     sessionRollResult.SessionRoll.SessionID,
		SessionRollID: sessionRollResult.SessionRoll.ID,
		PlayerID:      i.Member.User.ID,
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
			Content: fmt.Sprintf("Joined game with ID %s.", sessionRollID),
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

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Println("Failed to get interaction response:", err)
	}

	joinResult.SessionRoll.UpdatePlayerMsgID(i.Member.User.ID, msg.ID)

	// Update the Player with the message id
	_, err = c.manager.UpdateSessionRoll(context.Background(), &ronnied_actions.UpdateSessionRollInput{
		SessionRoll: joinResult.SessionRoll,
	})

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

	// Respond with the roll
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Content:    content,
			Components: i.Message.Components,
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

	embed := &discordgo.MessageEmbed{
		Title:       "Rollem Players Joining",
		Description: fmt.Sprintf("ID: %s", result.Session.ID),
		Color:       0x00ff00, // Green color
	}

	for _, participant := range sessionRollResult.SessionRoll.Players {
		user, userErr := s.User(participant.ID)
		if userErr != nil {
			log.Println("Failed to get user:", userErr)
			break
		}

		addDefault := true
		if entry := sessionRollResult.SessionRoll.HasPlayerEntry(participant.ID); entry != nil {
			addDefault = false
			user, userErr = s.User(entry.PlayerID)
			if userErr != nil {
				log.Println("Failed to get user:", userErr)
				continue
			}

			if entry.AssignedTo != "" {
				assignedUser, userErr := s.User(entry.AssignedTo)
				if userErr != nil {
					log.Println("Failed to get user:", userErr)
					continue
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   user.Username,
					Value:  fmt.Sprintf("Rolled a %d, assigned to %s", entry.Roll, assignedUser.Username),
					Inline: true,
				})
				addDefault = false
				continue
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   user.Username,
				Value:  fmt.Sprintf("Rolled a %d", entry.Roll),
				Inline: true,
			})
		}

		if addDefault {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   user.Username,
				Value:  "waiting for roll",
				Inline: true,
			})
		}
	}

	_, err = s.ChannelMessageEditEmbed(i.ChannelID, result.Session.MessageID, embed)
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
				Name:   "Host",
				Value:  i.Member.User.Username,
				Inline: true,
			},
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
		Label:    "Join Game",
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
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{joinButton, newButton}},
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
