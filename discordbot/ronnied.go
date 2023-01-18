package discordbot

import (
	"log"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (b *bot) addRonnieDCommand() error {
	log.Println("Adding command 'ronnied'...")

	cmd, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, b.guildID, &discordgo.ApplicationCommand{
		Name:        "ronnied",
		Description: "Ask Ronnie D for advice",
	})
	if err != nil {
		log.Printf("Error creating '/ronnied' command: %v", err)

		return err
	}

	b.registeredCommands = append(b.registeredCommands, cmd)

	return nil
}

func (b *bot) processRonnieD(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
