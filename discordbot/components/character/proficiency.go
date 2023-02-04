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
	"golang.org/x/sync/errgroup"
)

// Selecting proficiency options
func (c *Character) generateProficiencyChoices(char *entities.Character, choices []*entities.Choice) (string, []discordgo.MessageComponent, error) {
	if char.Class == nil {
		log.Println("Class is nil")
		return "", nil, errors.New("class is nil")
	}

	if len(choices) == 0 {
		log.Println("len(choices) == 0 ")
		return "", nil, dnderr.NewResourceExhaustedError("no proficiency choices")
	}

	var selectedChoice *entities.Choice
	for _, choice := range choices {
		if len(choice.Options) == 0 {
			log.Println("len(choice.Options) == 0")
			return "", nil, dnderr.NewResourceExhaustedError("no proficiency choices")
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

	err := c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency, choices)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}

	if selectedChoice == nil {
		log.Println("No proficiency choices exhausted")
		return "", nil, dnderr.NewResourceExhaustedError("no proficiency choices")
	}

	msg := fmt.Sprintf("Select %d starting proficiencies:", selectedChoice.Count)

	options := make([]discordgo.SelectMenuOption, len(selectedChoice.Options))

	for idx, choice := range selectedChoice.Options {
		if choice.GetOptionType() == entities.OptionTypeChoice {
			options[idx] = discordgo.SelectMenuOption{
				Label: choice.GetName(),
				Value: fmt.Sprintf("%s::%s::%d", choice.GetOptionType(), choice.GetKey(), idx),
			}
		} else {
			options[idx] = discordgo.SelectMenuOption{
				Label: choice.GetName(),
				Value: fmt.Sprintf("%s::%s::%s", choice.GetOptionType(), choice.GetKey(), choice.GetName()),
			}
		}

	}

	components := []discordgo.MessageComponent{
		discordgo.SelectMenu{
			MinValues: &selectedChoice.Count,
			MaxValues: selectedChoice.Count,
			CustomID:  selectProficiencyAction,
			Options:   options,
		},
	}

	return msg, components, nil
}

func (c *Character) handleProficiencySelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("handleProficiencySelect")

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepProficiency,
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

	choices, err := c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var done = true

	g, runCtx := errgroup.WithContext(context.Background())

	for _, choice := range choices {
		if choice.Status == entities.ChoiceStatusSelected {
			continue
		}

		if choice.Status == entities.ChoiceStatusActive {
			selectedChoiceIndex := -1

			// here is where I would reload the options for a given choice option
			for _, value := range i.MessageComponentData().Values {
				value := value

				g.Go(func() error {
					parts := strings.Split(value, "::")
					log.Println(parts)
					if parts[0] == string(entities.OptionTypeChoice) {
						// if we have a choice, we will check which was choses, set that to active and pass that choice back in
						// get index and iteract through options, setting other indexes to inactive and this to active, feed back into choice
						selectedChoiceIndex, err = strconv.Atoi(parts[2])
						if err != nil {
							log.Println(err)
							return err
						}
					} else {
						char, err = c.charManager.AddProficiency(runCtx, char, &entities.ReferenceItem{
							Key:  parts[1],
							Type: entities.ReferenceTypeProficiency,
						})
						if err != nil {
							log.Println(err)
							return err
						}
					}

					return nil
				})
			}

			err = g.Wait()
			if err != nil {
				log.Println(err)
				return // TODO handle error
			}

			choice.Status = entities.ChoiceStatusSelected

			// gross, but I need to get the last choice to set other options inactive
			if selectedChoiceIndex >= 0 {
				choice.Status = entities.ChoiceStatusActive
				for idx, option := range choice.Options {
					if option.GetOptionType() == entities.OptionTypeChoice {
						choiceOption := option.(*entities.Choice)
						if idx == selectedChoiceIndex {
							log.Println("choice: ", choiceOption.GetName(), " active")
							choiceOption.Status = entities.ChoiceStatusActive
						} else {
							log.Println("choice: ", choiceOption.GetName(), " inactive")
							choiceOption.Status = entities.ChoiceStatusInactive
						}
					}
				}
			}

			done = false
			break
		}
	}
	err = g.Wait()
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	err = c.charManager.SaveChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency, choices)
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
		log.Println(err)
		return // TODO handle error
	}

	log.Println("proficiency done: ", done)
	if done {
		c.handleEquipmentStep(s, i)
	} else {
		c.handleProficiencyStep(s, i)
	}
}

func (c *Character) handleProficiencyStep(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	state, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepProficiency,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	var proficiencyChoices []*entities.Choice

	if state.Step == entities.CreateStepRoll {
		proficiencyChoices = char.Class.ProficiencyChoices
		if char.Race != nil {
			if char.Race.StartingProficiencyOptions != nil {
				proficiencyChoices = append(proficiencyChoices, char.Race.StartingProficiencyOptions)
			}
		}
	} else {
		var choicesErr error
		proficiencyChoices, choicesErr = c.charManager.GetChoices(context.Background(), char.ID, entities.ChoiceTypeProficiency)
		if choicesErr != nil {
			log.Println(choicesErr)

			return // TODO handle error
		}
	}

	msg, components, err := c.generateProficiencyChoices(char, proficiencyChoices)
	if err != nil {
		log.Println(err)
		var exhaustedErr *dnderr.ResourceExhaustedError
		if errors.As(err, &exhaustedErr) {
			c.handleEquipmentStep(s, i)
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
