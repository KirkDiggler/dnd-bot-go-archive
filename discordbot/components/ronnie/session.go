package ronnie

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/ronnied_actions"
	"github.com/bwmarrin/discordgo"
	"log"
)

const socialGameID = "ronnie-rollem"

func updateSessionMessage(s *discordgo.Session, i *discordgo.InteractionCreate, sessionID string, participants []string) {
	// Create an updated embed with new participant details
	embed := &discordgo.MessageEmbed{
		Title:       "Game Session",
		Description: fmt.Sprintf("Session ID: %s", sessionID),
		Fields:      []*discordgo.MessageEmbedField{},
		Color:       0x00ff00,
	}

	for _, participant := range participants {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Participant",
			Value:  participant,
			Inline: true,
		})
	}

	// Create action buttons if needed
	actionButton := discordgo.Button{
		Label:    "Take Action",
		CustomID: "take_action_" + sessionID,
		Style:    discordgo.PrimaryButton,
	}

	// Update the message
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{Components: []discordgo.MessageComponent{actionButton}},
			},
		},
	})
	if err != nil {
		log.Println("Failed to update interaction:", err)
	}
}

func (c *RonnieD) SessionJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the session ID from the CustomID
	data := i.MessageComponentData()
	sessionID := data.CustomID[len("join_session_"):]
	log.Println("sessionID", sessionID)
	// Join the session
	_, err := c.manager.JoinSession(context.Background(), &ronnied_actions.JoinSessionInput{
		SessionID: sessionID,
		PlayerID:  i.Member.User.ID,
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
			Content: fmt.Sprintf("You have joined the session with ID %s.", sessionID),
		},
	})
	if err != nil {
		log.Println("Failed to respond to interaction:", err)
	}
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

	// Create an embed for the session
	embed := &discordgo.MessageEmbed{
		Title:       "New Game Session Created",
		Description: fmt.Sprintf("Session ID: %s", result.Session.ID),
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

	// Create a join button
	joinButton := discordgo.Button{
		Label:    "Join Session",
		CustomID: "join_session_" + result.Session.ID, // Include session ID in CustomID
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
}
