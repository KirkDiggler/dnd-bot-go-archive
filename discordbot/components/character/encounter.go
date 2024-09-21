package character

import (
	"context"
	"fmt"
	"log"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/bwmarrin/discordgo"
)

type rederPlayerCardInput struct {
	section string
	char    *entities.Character
}

func (c *Character) renderPlayerCard(s *discordgo.Session, i *discordgo.InteractionCreate, input *rederPlayerCardInput) {
	charEmbed := &discordgo.MessageEmbed{
		Description: input.char.NameString(),
	}

	char := input.char
	embeds := []*discordgo.MessageEmbed{charEmbed}

	embed := &discordgo.MessageEmbed{}
	switch input.section {
	case "stats":
		embed.Title = "Stats"

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Speed",
			Value:  fmt.Sprintf("%d ft", char.Speed),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "AC",
			Value:  fmt.Sprintf("%d", char.AC),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Max Hit Points",
			Value:  fmt.Sprintf("%d", char.MaxHitPoints),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Current Hit Points",
			Value:  fmt.Sprintf("%d", char.CurrentHitPoints),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Level",
			Value:  fmt.Sprintf("%d", char.Level),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Experience",
			Value:  fmt.Sprintf("%d", char.Experience),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Hit Die",
			Value:  fmt.Sprintf("%d", char.HitDie),
			Inline: true,
		})

		embeds = append(embeds, embed)
	case "attributes":
		embed.Title = "Attributes"
		for attr, score := range char.Attributes {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   string(attr),
				Value:  score.String(),
				Inline: true,
			})
		}

		embeds = append(embeds, embed)
	case "equipment":
		equippedEmbed := &discordgo.MessageEmbed{
			Title: "Equipped",
		}

		for slot, item := range char.EquippedSlots {
			if item == nil {
				log.Println("item is nil in slot", slot)
				continue
			}

			equippedEmbed.Fields = append(equippedEmbed.Fields, &discordgo.MessageEmbedField{
				Name:   string(slot),
				Value:  item.GetName(),
				Inline: true,
			})
		}

		embed.Title = "Backpack"
		for key := range char.Inventory {
			for _, item := range char.Inventory[key] {

				if !char.IsEquipped(item) {
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:   string(item.GetSlot()),
						Value:  item.GetName(),
						Inline: true,
					})
				}

			}

		}

		embeds = append(embeds, equippedEmbed)
		embeds = append(embeds, embed)

	case "proficiencies":

		embed.Title = "Proficiencies"
		embeds = append(embeds, embed)

		for _, key := range entities.ProficiencyTypes {
			if char.Proficiencies[key] == nil {
				continue
			}

			profEmbed := &discordgo.MessageEmbed{
				Title: string(key),
			}

			for _, prof := range char.Proficiencies[key] {
				profEmbed.Fields = append(profEmbed.Fields, &discordgo.MessageEmbedField{
					Value: prof.Name,
				})
			}

			embeds = append(embeds, profEmbed)
		}
	}

	buttonRow := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Stats",
					CustomID: "char:" + char.ID + ":stats",
					Style:    discordgo.PrimaryButton,
				},
				discordgo.Button{
					Label:    "Attributes",
					CustomID: "char:" + char.ID + ":attributes",
					Style:    discordgo.PrimaryButton,
				},
				discordgo.Button{
					Label:    "Equipment",
					CustomID: "char:" + char.ID + ":equipment",
					Style:    discordgo.PrimaryButton,
				},
				discordgo.Button{
					Label:    "Proficiencies",
					CustomID: "char:" + char.ID + ":proficiencies",
					Style:    discordgo.PrimaryButton,
				},
			}},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Content:    "Here is your character",
			Embeds:     embeds,
			Components: buttonRow,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func (c *Character) handleEncounterJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// First make sure the user has a character created
	// if not, send them to the character creation flow
	// if they do, add them to the encounter
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	data := i.MessageComponentData()
	encounterID := data.CustomID[len("encounter:join:"):]
	encounter, err := c.charManager.GetEncounter(context.Background(), encounterID)
	if err != nil {
		log.Println(err)
		return
	}

	encounter.Players = append(encounter.Players, i.Member.User.ID)
	_, err = c.charManager.UpdateEncounter(context.Background(), encounter)
	if err != nil {
		log.Println(err)
		return
	}

	c.renderPlayerCard(s, i, &rederPlayerCardInput{
		section: "stats",
		char:    char,
	})
}

func (c *Character) handleShowStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	c.renderPlayerCard(s, i, &rederPlayerCardInput{
		section: "stats",
		char:    char,
	})
}

func (c *Character) handleShowProficiencies(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	c.renderPlayerCard(s, i, &rederPlayerCardInput{
		section: "proficiencies",
		char:    char,
	})
}

func (c *Character) handleShowEquipment(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	c.renderPlayerCard(s, i, &rederPlayerCardInput{
		section: "equipment",
		char:    char,
	})
}

func (c *Character) handleShowAttributes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	c.renderPlayerCard(s, i, &rederPlayerCardInput{
		section: "attributes",
		char:    char,
	})
}

func (c *Character) handleEncounterCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	encounter, err := c.charManager.CreateEncounter(context.Background(), &entities.Encounter{
		Players: []string{},
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Encounter",
		Description: "This is the encounter",
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Join Encounter",
							Style:    discordgo.SuccessButton,
							CustomID: fmt.Sprintf("encounter:join:%s", encounter.ID),
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// grabe the response and set the message id
	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Println(err)
		return
	}

	// save the message id to the encounter
	encounter.MessageID = msg.ID
	_, err = c.charManager.UpdateEncounter(context.Background(), encounter)
	if err != nil {
		log.Println(err)
		return
	}
}
