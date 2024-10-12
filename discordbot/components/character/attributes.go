package character

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/bwmarrin/discordgo"
)

func rollAttributes() ([]*dice.RollResult, error) {
	attributes := make([]*dice.RollResult, 6)

	for idx := range attributes {
		roll, err := dice.RollString("4d6")
		if err != nil {
			return nil, err
		}
		attributes[idx] = roll
	}
	return attributes, nil
}

// Setting Attributes
func (c *Character) handleAttributeSelect(s *discordgo.Session, i *discordgo.InteractionCreate, attribute string, selectSlice []string) {
	draft, err := c.charManager.GetDraft(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	char, err := c.charManager.Get(context.Background(), draft.CharacterID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}
	idx, err := strconv.Atoi(selectSlice[1])
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	if idx >= len(char.Rolls) {
		log.Printf("index out of range. idx: %d, len: %d", idx, len(char.Rolls))
		return // TODO: Handle error
	}
	// TODO: make set attribut function that returns bool if it was set
	if !char.Rolls[idx].Used { // We have not used this one
		char.AddAttribute(entities.Attribute(attribute), char.Rolls[idx].Total-char.Rolls[idx].Lowest)
		log.Printf("setting %s to %s ", attribute, char.Attributes[entities.Attribute(attribute)])
		char.Rolls[idx].Used = true
		// TODO Calculate modifiers
	}

	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	_, err = c.charManager.UpdateDraft(context.Background(), draft, char)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	rolls := char.Rolls
	attributeSelectData, err := c.generateAttributeSelect(char, rolls, i)
	done := false
	if err != nil {
		if err.Error() == "done" {
			done = true
		} else {
			log.Println(err)
			return // TODO: Handle error
		}
	}
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("Rolls: ")
	for _, roll := range rolls {
		msgBuilder.WriteString(fmt.Sprintf("%d, ", roll.Total-roll.Lowest))
	}

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: char.ID,
		OwnerID:     i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	oldInteraction := &discordgo.Interaction{AppID: i.AppID, Token: state.LastToken}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	var response *discordgo.InteractionResponse
	if done {
		c.handleProficiencyStep(s, i)
		return
	} else {
		response = &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:      discordgo.MessageFlagsEphemeral,
				Content:    msgBuilder.String(),
				Components: attributeSelectData,
			},
		}
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
}

func (c *Character) generateAttributeSelect(char *entities.Character, rolls []*dice.RollResult, i *discordgo.InteractionCreate) ([]discordgo.MessageComponent, error) {
	userID := i.Member.User.ID

	selectionOrder := []string{"Str", "Dex", "Con", "Int", "Wis", "Cha"}

	selected := make(map[entities.Attribute]*entities.AbilityScore)

	for k, v := range char.Attributes {
		selected[k] = v
	}

	attrToSelect := ""
	for _, attr := range selectionOrder {
		if selected[entities.Attribute(attr)].Score == 0 {
			attrToSelect = attr
			break
		}
	}

	if attrToSelect == "" {
		return nil, errors.New("done")
	}

	components := make([]discordgo.SelectMenuOption, 0)
	for idx, roll := range rolls {
		log.Println("roll: ", roll, "used: ", roll.Used)
		if !roll.Used {
			components = append(components, discordgo.SelectMenuOption{
				Label: fmt.Sprintf("%v  %d", roll.Rolls, roll.Total-roll.Lowest),
				Value: fmt.Sprintf("roll:%d:%d", idx, roll.Total-roll.Lowest),
			})
		}
	}

	response := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "STR", selected[entities.AttributeStrength].Score),
					Disabled: selected[entities.AttributeStrength].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:str", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "DEX", selected[entities.AttributeDexterity].Score),
					Disabled: selected[entities.AttributeDexterity].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:dex", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "CON", selected[entities.AttributeConstitution].Score),
					Disabled: selected[entities.AttributeConstitution].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:con", buttonAttributeKey, userID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "INT", selected[entities.AttributeIntelligence].Score),
					Disabled: selected[entities.AttributeIntelligence].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:int", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "WIS", selected[entities.AttributeWisdom].Score),
					Disabled: selected[entities.AttributeWisdom].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:wis", buttonAttributeKey, userID),
				},
				discordgo.Button{
					Label:    fmt.Sprintf("%s: %d", "CHA", selected[entities.AttributeCharisma].Score),
					Disabled: selected[entities.AttributeCharisma].Score == 0,
					CustomID: fmt.Sprintf("%s:%s:cha", buttonAttributeKey, userID),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					Placeholder: strings.ToUpper(attrToSelect),
					CustomID:    fmt.Sprintf("%s:%s:%s", selectAttributeKey, userID, attrToSelect),
					Options:     components,
				},
			},
		},
	}

	return response, nil
}

func (c *Character) handleRollCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	draft, err := c.charManager.GetDraft(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	char, err := c.charManager.Get(context.Background(), draft.CharacterID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	log.Println("char: ", char.String())

	log.Println("Rolling for", char.Name, "the ", char.Race.Key, " ", char.Class.Key)

	_, err = c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: char.ID,
		OwnerID:     i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	rolls, err := rollAttributes()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	char.Rolls = rolls
	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println("error returned from charManager.Put: ", err)
		return // TODO: Handle error
	}

	attributeSelectData, err := c.generateAttributeSelect(char, rolls, i)
	if err != nil {
		log.Println("error returned from generateAttributeSelect: ", err)
		return // TODO handle error
	}
	msgBuilder := strings.Builder{}
	msgBuilder.WriteString("Rolls: ")
	for _, roll := range rolls {
		msgBuilder.WriteString(fmt.Sprintf("%d, ", roll.Total-roll.Lowest))
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Content:    msgBuilder.String(),
			Components: attributeSelectData,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println("error returned from InteractionRespond: ", err)
		return // TODO handle error
	}
}

// Go through choice