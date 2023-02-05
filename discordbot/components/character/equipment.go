package character

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/bwmarrin/discordgo"
)

// Selecting equipment options
func (c *Character) generateStartingEquipmentChoices(char *entities.Character, choices []*entities.Choice) (string, []discordgo.MessageComponent, error) {
	if char.Class == nil {
		log.Println("Class is nil")
		return "", nil, errors.New("class is nil")
	}

	if len(choices) == 0 {
		log.Println("No equipment choices")
		return "", nil, errors.New("equipment len(choices) == 0")
	}

	var selectedChoice *entities.Choice

	for _, choice := range choices {
		if len(choice.Options) == 0 {
			log.Println("equipment len(choice.Options) == 0")
			return "", nil, errors.New("no equipment choices")
		}

		var SelectErr error
		selectedChoice, SelectErr = c.getNextChoiceOption(choice)
		if SelectErr != nil {
			var exhaustedErr *dnderr.ResourceExhaustedError
			if errors.As(SelectErr, &exhaustedErr) {
				continue
			}

			log.Println(SelectErr)
			return "", nil, SelectErr
		}

		selectedChoice.Status = entities.ChoiceStatusActive

		break
	}

	err := c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeEquipment, choices)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	if selectedChoice == nil {
		log.Println("No equipment choices")
		return "", nil, dnderr.NewResourceExhaustedError("no equipment choices")
	}

	msg := fmt.Sprintf("Select %d starting equipment:", selectedChoice.Count)

	options := make([]discordgo.SelectMenuOption, len(selectedChoice.Options))

	log.Println("selectedChoice.Options count: ", len(selectedChoice.Options))
	for idx, choice := range selectedChoice.Options {
		if choice.GetOptionType() == entities.OptionTypeChoice {
			options[idx] = discordgo.SelectMenuOption{
				Label: choice.GetName(),
				Value: fmt.Sprintf("%s:%s:%d", choice.GetOptionType(), choice.GetKey(), idx),
			}
		} else {
			options[idx] = discordgo.SelectMenuOption{
				Label: choice.GetName(),
				Value: fmt.Sprintf("%s:%s", choice.GetOptionType(), choice.GetKey()),
			}
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.SelectMenu{
			MinValues: &selectedChoice.Count,
			MaxValues: selectedChoice.Count,
			CustomID:  selectEquipmentAction,
			Options:   options,
		},
	}

	return msg, components, nil
}

func (c *Character) handleEquipmentSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepEquipment,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	choices, err := c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeEquipment)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var done = true

	for _, choice := range choices {
		if choice.Status == entities.ChoiceStatusSelected {
			continue
		}

		if choice.Status == entities.ChoiceStatusActive {
			done = false
			selectedChoiceIndex := -1

			// here is where I would reload the options for a given choice option
			for _, value := range i.MessageComponentData().Values {
				parts := strings.Split(value, ":")
				if parts[0] == string(entities.OptionTypeChoice) {
					// get index and iteract through options, setting other indexes to inactive and this to active, feed back into choice
					selectedChoiceIndex, err = strconv.Atoi(parts[2])
					if err != nil {
						log.Println(err)
						return // TODO handle error
					}
				} else {
					log.Println("add equipment", parts)
					// add equipment
					char, err = c.charManager.AddInventory(context.Background(), char, parts[1])
					if err != nil {
						log.Println(err)
						return // TODO handle error
					}
				}
			}
			choice.Status = entities.ChoiceStatusSelected

			// gross, but I need to get the last choice to set other options inactive
			if selectedChoiceIndex >= 0 {
				choice.Status = entities.ChoiceStatusActive
				for idx, option := range choice.Options {
					if idx == selectedChoiceIndex {
						log.Println("choice: ", option.GetName(), " active")
						option.SetStatus(entities.ChoiceStatusActive)
					} else {
						log.Println("choice: ", option.GetName(), " inactive")
						option.SetStatus(entities.ChoiceStatusInactive)
					}
				}
			}

			break
		}

		log.Println("done choosing equipment")

		done = true
	}

	err = c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeEquipment, choices)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	oldInteraction := &discordgo.Interaction{
		AppID: i.AppID,
		Token: state.LastToken,
	}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println("error with InteractionResponseDelete: ", err)
		return // TODO handle error
	}

	if done {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "equipment selected",
			},
		}

		err = s.InteractionRespond(i.Interaction, response)
		if err != nil {
			log.Println(err)
			return // TODO handle error
		}
	} else {
		c.handleEquipmentStep(s, i)
	}
}

func (c *Character) handleEquipmentStep(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepEquipment,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var startingEquipementChoices []*entities.Choice

	if state.Step == entities.CreateStepProficiency {
		log.Println("getting starting equipment choices ", char.Class.StartingEquipmentChoices)
		startingEquipementChoices = char.Class.StartingEquipmentChoices
	} else {
		var choicesErr error
		if choicesErr != nil {
			var notFoundErr *dnderr.NotFoundError
			if !errors.As(choicesErr, &notFoundErr) {
				log.Println(choicesErr)
				return // TODO handle error
			}
		}
		startingEquipementChoices, choicesErr = c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeEquipment)
		if choicesErr != nil {
			log.Println(choicesErr)
			return // TODO handle error

		}
	}

	msg, components, err := c.generateStartingEquipmentChoices(char, startingEquipementChoices)
	if err != nil {
		log.Println(err)
		var exhaustedErr *dnderr.ResourceExhaustedError
		if errors.As(err, &exhaustedErr) {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "equipment selected",
				},
			}

			err = s.InteractionRespond(i.Interaction, response)
			if err != nil {
				log.Println(err)
				return // TODO handle error
			}
		}

		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: msg,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}
