package character

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (c *Character) handleAttack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
	attack, err := char.Attack()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
	msgBuilder := strings.Builder{}

	for _, a := range attack {
		msgBuilder.WriteString(a.String())
		msgBuilder.WriteString("\n")
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msgBuilder.String(),
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}
