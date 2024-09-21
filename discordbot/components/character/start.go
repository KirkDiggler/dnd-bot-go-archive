package character

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/errgroup"
)

func (c *Character) startNewChoices(number int) ([]*charChoice, error) {
	// TODO cache these. cache in the client wrapper? configurable?
	races, err := c.client.ListRaces()
	if err != nil {
		return nil, err
	}

	classes, err := c.client.ListClasses()
	if err != nil {
		return nil, err
	}
	log.Println("Starting new choices")
	choices := make([]*charChoice, number)

	var monk *entities.Class
	for _, class := range classes {
		if class.Key == "fighter" {
			monk = class
			break
		}
	}
	for idx := 0; idx < number-1; idx++ {
		class := classes[rand.Intn(len(classes))]
		time.Sleep(1 * time.Nanosecond)
		race := races[rand.Intn(len(races))]
		choices[idx] = &charChoice{
			Race:  race,
			Class: class,
		}
	}
	choices[number-1] = &charChoice{
		Race:  races[rand.Intn(len(races))],
		Class: monk,
	}

	return choices, nil
}

func (c *Character) handleListCharacters(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Handling list characters")
	characters, err := c.charManager.List(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	// Create select menu options from the list of characters
	options := make([]discordgo.SelectMenuOption, len(characters))
	for idx, char := range characters {
		options[idx] = discordgo.SelectMenuOption{
			Label: char.Name,
			Value: char.ID,
		}
	}

	// Create the response with the select menu
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Select a character:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "select_character",
							Placeholder: "Choose a character",
							Options:     options,
						},
					},
				},
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
	}
}

func (c *Character) handleNewCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Handling new character")

	char, err := c.charManager.Put(context.Background(), &entities.Character{
		OwnerID: i.Member.User.ID,
	})
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	initialState := &entities.CharacterCreation{
		CharacterID: char.ID, // Use the generated UUID from the response
		OwnerID:     i.Member.User.ID, // Set the OwnerID
		LastToken:   i.Token,
	}
	_, err = c.charManager.SaveState(context.Background(), initialState)
	if err != nil {
		log.Println(err)
	}

	err = c.renderCharacterCreate(s, i, initialState)
	if err != nil {
		log.Println(err)
	}
}

func (c *Character) createRaceOptions() []discordgo.SelectMenuOption {
	races, err := c.client.ListRaces()
	if err != nil {
		log.Println(err)
		return make([]discordgo.SelectMenuOption, 0)
	}
	raceOptions := make([]discordgo.SelectMenuOption, len(races))
	for idx, race := range races {
		raceOptions[idx] = discordgo.SelectMenuOption{
			Label: race.Name,
			Value: race.Key,
		}
	}

	return raceOptions
}

func (c *Character) createClassOptions() []discordgo.SelectMenuOption {
	classes, err := c.client.ListClasses()
	if err != nil {
		log.Println(err)
		return make([]discordgo.SelectMenuOption, 0)
	}
	classOptions := make([]discordgo.SelectMenuOption, len(classes))
	for idx, class := range classes {
		classOptions[idx] = discordgo.SelectMenuOption{
			Label: class.Name,
			Value: class.Key,
		}
	}

	return classOptions
}

// Selecting a character
func (c *Character) handleRandomStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	choices, err := c.startNewChoices(4)
	if err != nil {
		log.Println(err)
		// TODO: Handle error
		return
	}

	_, err = c.charManager.SaveState(context.Background(), &entities.CharacterCreation{
		CharacterID: i.Member.User.ID,
		LastToken:   i.Token,
		Step:        entities.CreateStepSelect,
	})
	if err != nil {
		log.Println(err)
	}

	components := make([]discordgo.SelectMenuOption, len(choices))
	for idx, char := range choices {
		components[idx] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%s the %s %s", i.Member.User.Username, char.Race.Name, char.Class.Name),
			Value: fmt.Sprintf("choice:%d:%s:%s", idx, char.Race.Key, char.Class.Key),
		}
	}
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Select your new character:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							// Select menu, as other components, must have a customID, so we set it to this value.
							CustomID:    selectCaracterAction,
							Placeholder: "This is the tale of...👇",
							Options:     components,
						},
					},
				},
			},
		},
	}
	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Character) renderCharacterCreate(s *discordgo.Session, i *discordgo.InteractionCreate, state *entities.CharacterCreation) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Create your character:",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    selectRaceAction,
							Placeholder: "Select your race",
							Options:     c.createRaceOptions(),
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    selectClassAction,
							Placeholder: "Select your class",
							Options:     c.createClassOptions(),
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: submitCharacterStart,
							Label:    "Submit",
							Style:    discordgo.PrimaryButton,
						},
					},
				},
			},
		},
	}

	return s.InteractionRespond(i.Interaction, response)
}

func handleSubmitNewCharacterInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	modal := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "character_name_modal",
			Title:    "Enter Character Name",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "input_name",
							Style:       discordgo.TextInputShort,
							Label:       "Character Name",
							Placeholder: "Enter your character's name",
							Required:    true,
						},
					},
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, modal)
	if err != nil {
		log.Println(err)

		return
	}
}

func (c *Character) handleNameCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()
	log.Println("Handling name character", data)
	name := ""
	for _, component := range data.Components {
		if actionRow, ok := component.(*discordgo.ActionsRow); ok {
			for _, comp := range actionRow.Components {
				if input, inputOK := comp.(*discordgo.TextInput); inputOK && input.CustomID == "input_name" {
					name = input.Value
					break
				}
			}
			if name != "" {
				break
			}
		}
	}

	if name == "" {
		log.Println("Error: No name input found")
		return
	}

	state, err := c.charManager.GetState(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	char, err := c.charManager.Get(context.Background(), state.CharacterID)
	if err != nil {
		log.Println(err)
		return
	}

	char.Name = name
	_, err = c.charManager.Put(context.Background(), char)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Character Name", char.Name)

	oldInteraction := &discordgo.Interaction{
		AppID: i.AppID,
		Token: state.LastToken,
	}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}
	c.handleRollCharacter(s, i)
}

func (c *Character) initializeCharacter(charID string) error {
	log.Println("Initializing character", charID)
	char, err := c.charManager.Get(context.Background(), charID)
	if err != nil {
		return err
	}

	if char.Race == nil {
		return dnderr.NewInvalidEntityError("Race is nil")
	}

	if char.Class == nil {
		return dnderr.NewInvalidEntityError("Class is nil")
	}
	g, runCtx := errgroup.WithContext(context.Background())

	char.HitDie = char.Class.HitDie
	char.AC = 10
	char.Level = 1

	char.Speed = char.Race.Speed

	// Load the race starting data
	for _, prof := range char.Race.StartingProficiencies {
		prof := prof
		g.Go(func() error {
			char, err = c.charManager.AddProficiency(runCtx, char, prof)
			if err != nil {
				return err
			}

			return nil
		})
	}

	for _, bonus := range char.Race.AbilityBonuses {
		char.AddAbilityBonus(bonus)
	}

	for _, prof := range char.Class.Proficiencies {
		prof := prof
		g.Go(func() error {
			char, err = c.charManager.AddProficiency(runCtx, char, prof)
			if err != nil {
				return err
			}

			return nil
		})
	}

	for _, equip := range char.Class.StartingEquipment {
		equip := equip
		g.Go(func() error {
			char, err = c.charManager.AddInventory(runCtx, char, equip.Equipment.Key)
			if err != nil {
				return err
			}

			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		return err
	}
	_, err = c.charManager.Put(context.Background(), char)

	log.Println("Character initialized", charID)

	return err
}

func (c *Character) handleRaceAndClassSelection(s *discordgo.Session, i *discordgo.InteractionCreate) {
	state, err := c.charManager.GetState(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	existing, err := c.charManager.Get(context.Background(), state.CharacterID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	data := i.MessageComponentData()
	selection := data.Values[0] // Assuming single selection

	log.Println("Selection", selection)

	// Store the selection based on the CustomID
	switch data.CustomID {
	case "select-race":
		state.Steps |= entities.SelectRaceStep
		existing.Race = &entities.Race{Key: selection}
	case "select-class":
		state.Steps |= entities.SelectClassStep
		existing.Class = &entities.Class{Key: selection}
	}

	char, err := c.charManager.Put(context.Background(), existing)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	_, err = c.charManager.SaveState(context.Background(), state)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	_, err = c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: state.CharacterID, // Use the character ID from the state
		OwnerID:     state.OwnerID,     // Ensure OwnerID is maintained
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if state.Steps&entities.SelectRaceStep != 0 && state.Steps&entities.SelectClassStep != 0 {
		err = c.initializeCharacter(char.ID)
		if err != nil {
			log.Println(err)
			return // TODO handle error
		}
	}

	// Acknowledge the interaction to prevent timeout
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Character) handleCharSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Handling character select")
	selectString := strings.Split(i.MessageComponentData().Values[0], ":")
	if len(selectString) != 4 {
		log.Printf("Invalid select string: %s", selectString)
		return
	}

	race := selectString[2]
	class := selectString[3]

	state, err := c.charManager.GetState(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO: Handle error
	}

	char, err := c.charManager.Put(context.Background(), &entities.Character{
		ID:      state.CharacterID, // Use the character ID from the state
		OwnerID: state.OwnerID,     // Ensure OwnerID is maintained
		Name:    i.Member.User.Username,
		Race: &entities.Race{
			Key: race,
		},
		Class: &entities.Class{
			Key: class,
		},
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	err = c.initializeCharacter(char.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	lastState, err := c.getAndUpdateState(&entities.CharacterCreation{
		CharacterID: state.CharacterID, // Use the character ID from the state
		OwnerID:     state.OwnerID,     // Ensure OwnerID is maintained
		LastToken:   i.Token,
		Step:        entities.CreateStepRoll,
	})
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	oldInteraction := &discordgo.Interaction{
		AppID: i.AppID,
		Token: lastState.LastToken,
	}
	err = s.InteractionResponseDelete(oldInteraction)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	log.Println("Character selected", char.ID)

	c.handleRollCharacter(s, i)
}
