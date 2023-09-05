package character

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/internal/managers/rooms"

	"github.com/bwmarrin/discordgo"
)

func (c *Character) handleAttack(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roomResult, err := c.roomManager.LoadRoom(context.Background(), &rooms.LoadRoomInput{
		PlayerID: i.Member.User.ID,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if roomResult.Room.IsEmpty() {
		c.practiceAttack(s, i)

		return
	}

	msg, err := c.roomManager.Attack(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Character) practiceAttack(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
