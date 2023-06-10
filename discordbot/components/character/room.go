package character

import (
	"context"
	"fmt"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/rooms"
	"github.com/bwmarrin/discordgo"
	"log"
)

func (c *Character) handleLoadRoom(s *discordgo.Session, i *discordgo.InteractionCreate) {
	room, err := c.roomManager.LoadRoom(context.Background(), &rooms.LoadRoomInput{PlayerID: i.Member.User.ID})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: room.Room.String(),
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}
