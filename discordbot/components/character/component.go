package character

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/KirkDiggler/dnd-bot-go/internal/managers/characters"
	"github.com/KirkDiggler/dnd-bot-go/internal/managers/session"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/bwmarrin/discordgo"
)

const (
	selectRaceAction        = "select-race"
	selectClassAction       = "select-class"
	setNameAction           = "set-name"
	submitCharacterStart    = "submit-character-start"
	selectCaracterAction    = "select-character"
	equipInventoryAction    = "equip-inventory"
	selectProficiencyAction = "select-proficiency"
	selectEquipmentAction   = "select-equipment"
	rollCharacterAction     = "roll-character"
	selectAttributeKey      = "select-attribute"
	buttonAttributeKey      = "button-attribute"
)

type Character struct {
	client         dnd5e.Client
	charManager    characters.Manager
	sessionManager session.Manager
}

type CharacterConfig struct {
	Client           dnd5e.Client
	CharacterManager characters.Manager
	SessionManager   session.Manager
}

type charChoice struct {
	Name  string
	Race  *entities.Race
	Class *entities.Class
}

func NewCharacter(cfg *CharacterConfig) (*Character, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.CharacterManager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.CharacterManager")
	}

	if cfg.SessionManager == nil {
		return nil, dnderr.NewMissingParameterError("cfg.SessionManager")
	}

	return &Character{
		client:         cfg.Client,
		charManager:    cfg.CharacterManager,
		sessionManager: cfg.SessionManager,
	}, nil
}

func (c *Character) GetApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "character",
		Description: "Generate a character",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "new",
				Description: "Put a new character from a random list",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "load",
				Description: "Load your character",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "display",
				Description: "Display your character",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "equip",
				Description: "Equip an item from your inventory",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "attack",
				Description: "Attack a target using your equipped weapon",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			}, {
				Name:        "encounter",
				Description: "Start an encounter to join with other players",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	}
}

func (c *Character) GetDNDApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "dnd",
		Description: "DnD commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "character",
				Description: "Character related commands",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "create",
						Description: "Create a new character",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
					},
				},
			},
		},
	}
}

func (c *Character) HandleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case "character":
			switch i.ApplicationCommandData().Options[0].Name {
			case "random":
				c.handleRandomStart(s, i)
			case "load":
				c.handleLoadCharacter(s, i)
			case "display":
				c.handleDisplayCharacter(s, i)
			case "equip":
				c.handleEquipInventory(s, i)
			case "attack":
				c.handleAttack(s, i)
			case "encounter":
				c.handleEncounterCreate(s, i)
			}
		case "dnd":
			switch i.ApplicationCommandData().Options[0].Name {
			case "character":
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "create":
					log.Println("create character")
					c.handleNewCharacter(s, i)
				}
			case "list":
				c.handleListCharacters(s, i)
			}
		}
	case discordgo.InteractionModalSubmit:
		if i.ModalSubmitData().CustomID == "character_name_modal" {
			c.handleNameCharacter(s, i)
		}

	case discordgo.InteractionMessageComponent:
		strKey := fmt.Sprintf("%s:%s:Str", selectAttributeKey, i.Member.User.ID)
		dexKey := fmt.Sprintf("%s:%s:Dex", selectAttributeKey, i.Member.User.ID)
		conKey := fmt.Sprintf("%s:%s:Con", selectAttributeKey, i.Member.User.ID)
		intKey := fmt.Sprintf("%s:%s:Int", selectAttributeKey, i.Member.User.ID)
		wisKey := fmt.Sprintf("%s:%s:Wis", selectAttributeKey, i.Member.User.ID)
		chaKey := fmt.Sprintf("%s:%s:Cha", selectAttributeKey, i.Member.User.ID)

		switch i.MessageComponentData().CustomID {
		case selectRaceAction:
			log.Println("selectRaceAction")
			c.handleRaceAndClassSelection(s, i)
		case selectClassAction:
			log.Println("selectClassAction")
			c.handleRaceAndClassSelection(s, i)
		case setNameAction:
			log.Println("setNameAction")
			c.handleNameCharacter(s, i)
		case submitCharacterStart:
			log.Println("submitCharacterStart")
			handleSubmitNewCharacterInteraction(s, i)
		case selectCaracterAction:
			c.handleCharSelect(s, i)
		case rollCharacterAction:
			c.handleRollCharacter(s, i)
		case strKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Str", selectSlice)
		case dexKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Dex", selectSlice)
		case conKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Con", selectSlice)
		case intKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Int", selectSlice)
		case wisKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Wis", selectSlice)
		case chaKey:
			selectSlice := strings.Split(i.MessageComponentData().Values[0], ":")
			c.handleAttributeSelect(s, i, "Cha", selectSlice)
		case selectProficiencyAction:
			c.handleProficiencySelect(s, i)
		case selectEquipmentAction:
			c.handleEquipmentSelect(s, i)
		case equipInventoryAction:
			c.handleEquipInventorySelect(s, i)
		default:
			data := i.MessageComponentData()
			if strings.HasPrefix(data.CustomID, "encounter:join:") {
				c.handleEncounterJoin(s, i)
			}

			if strings.HasPrefix(data.CustomID, "char:") {
				if strings.HasSuffix(data.CustomID, ":stats") {
					c.handleShowStats(s, i)
				}

				if strings.HasSuffix(data.CustomID, ":attributes") {
					c.handleShowAttributes(s, i)
				}

				if strings.HasSuffix(data.CustomID, ":equipment") {
					c.handleShowEquipment(s, i)
				}

				if strings.HasSuffix(data.CustomID, ":proficiencies") {
					c.handleShowProficiencies(s, i)
				}
			}
		}
	}
}

// Go through choices, searching the active path and return the first unset options
func (c *Character) getNextChoiceOption(input *entities.Choice) (*entities.Choice, error) {
	if input == nil {
		return nil, dnderr.NewMissingParameterError("input")
	}

	switch input.Status {
	case entities.ChoiceStatusUnset:
		return input, nil
	case entities.ChoiceStatusActive:
		for _, option := range input.Options {
			if option.GetStatus() == entities.ChoiceStatusInactive {
				continue
			}

			if option.GetOptionType() == entities.OptionTypeChoice {
				optionChoice := option.(*entities.Choice)

				return c.getNextChoiceOption(optionChoice)
			}

			return input, nil
		}
	}
	return nil, dnderr.NewResourceExhaustedError("no active choice")
}

func (c *Character) handleDisplayCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: char.String(),
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return
	}

	// Update the last token in the session
	err = c.sessionManager.UpdateLastToken(context.Background(), i.Member.User.ID, i.Token)
	if err != nil {
		log.Println("Failed to update last token:", err)
		// Note: We don't return here as the main interaction was successful
	}
}

func (c *Character) handleLoadCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if char.Race == nil || char.Class == nil {
		log.Println("Character not fully loaded")
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Loaded character %s the %s %s", char.Name, char.Race.Name, char.Class.Name),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return
	}

	// Update the last token in the session
	err = c.sessionManager.UpdateLastToken(context.Background(), i.Member.User.ID, i.Token)
	if err != nil {
		log.Println("Failed to update last token:", err)
		// Note: We don't return here as the main interaction was successful
	}
}

func (c *Character) handleLoadCharacter(s *discordgo.Session, i *discordgo.InteractionCreate) {
	char, err := c.charManager.Get(context.Background(), i.Member.User.ID)
	if err != nil {
		log.Println(err)
		return // TODO handle error
	}

	if char.Race == nil || char.Class == nil {
		log.Println("Character not fully loaded")
		return // TODO handle error
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Loaded character %s the %s %s", char.Name, char.Race.Name, char.Class.Name),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	err = s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println(err)
		return
	}

	// Update the last token in the session
	err = c.sessionManager.UpdateLastToken(context.Background(), i.Member.User.ID, i.Token)
	if err != nil {
		log.Println("Failed to update last token:", err)
		// Note: We don't return here as the main interaction was successful
	}
}

session, err := c.sessionManager.GetWithDraft(context.Background(), i.Member.User.ID)
if err != nil {
	log.Println(err)
	return // TODO: Handle error
}

oldInteraction := &discordgo.Interaction{
	AppID: i.AppID,
	Token: session.LastToken,
}
err = s.InteractionResponseDelete(oldInteraction)
if err != nil {
	log.Println("Error deleting old interaction:", err)
	// Handle error as appropriate
}
