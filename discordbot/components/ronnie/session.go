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
	_, err = c.manager.JoinSession(context.Background(), &ronnied_actions.JoinSessionInput{
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
			Content: fmt.Sprintf("You have joined the session with ID %s.", sessionRollID),
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Rollem",
						CustomID: "rollem:" + sessionRollID,
						Style:    discordgo.PrimaryButton,
					},
				},
				},
			},
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
	if !sessionRollResult.SessionRoll.HasPlayer(i.Member.User.ID) {
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

	rollResult, err := c.manager.AddSessionRoll(context.Background(), &ronnied_actions.AddSessionRollInput{
		SessionRollID: sessionRollResult.SessionRoll.ID,
		PlayerID:      i.Member.User.ID,
	})
	if err != nil {
		fmt.Println("c.manager.AddSessionRoll returned err:", err)
		return
	}

	// Respond with the roll
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Content:    fmt.Sprintf("You rolled a %d.", rollResult.SessionEntry.Roll),
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
		Title:       "Game Session Players Joining",
		Description: fmt.Sprintf("Session ID: %s", result.Session.ID),
		Color:       0x00ff00, // Green color
	}

	for _, participant := range sessionRollResult.SessionRoll.Players {
		user, userErr := s.User(participant)
		if userErr != nil {
			log.Println("Failed to get user:", userErr)
			continue
		}

		addDefault := true
		for _, entry := range sessionRollResult.SessionRoll.Entries {
			if entry.PlayerID == participant {
				user, userErr = s.User(entry.PlayerID)
				if userErr != nil {
					log.Println("Failed to get user:", userErr)
					continue
				}

				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   user.Username,
					Value:  fmt.Sprintf("Rolled a %d", entry.Roll),
					Inline: true,
				})
				addDefault = false
				break
			}
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
		Title:       "New Game Session Created",
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
		user, userErr := s.User(participant)
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
		Label:    "Join Session",
		CustomID: "join_session:" + rollResult.SessionRoll.ID, // Include session ID in CustomID
		Style:    discordgo.PrimaryButton,
	}

	// Respond with the embed and button
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{joinButton}},
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
