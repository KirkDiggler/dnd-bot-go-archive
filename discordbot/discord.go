package discordbot

import (
	"log"

	"github.com/KirkDiggler/dnd-bot-go/clients/dnd5e"
	"github.com/KirkDiggler/dnd-bot-go/discordbot/components/character"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session *discordgo.Session
	guildID string

	registeredCommands []*discordgo.ApplicationCommand
	characterComponent *character.Component
}

type Config struct {
	Token   string
	GuildID string
	AppID   string
	Client  dnd5e.Interface
}

func New(cfg *Config) (*bot, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	if cfg.Client == nil {
		return nil, dnderr.NewMissingParameterError("cfg.Client")
	}

	if cfg.Token == "" {
		return nil, dnderr.NewMissingParameterError("cfg.Token")
	}

	if cfg.GuildID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.GuildID")
	}

	if cfg.AppID == "" {
		return nil, dnderr.NewMissingParameterError("cfg.AppID")
	}
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents |= discordgo.IntentGuildMembers
	session.Identify.Intents |= discordgo.IntentsGuilds
	session.Identify.Intents |= discordgo.IntentsGuildMessageReactions
	component, err := character.New(&character.Config{
		Client:  cfg.Client,
		Session: session,
		AppID:   cfg.AppID,
		GuildID: cfg.GuildID,
	})
	if err != nil {
		return nil, err
	}
	return &bot{
		session:            session,
		guildID:            cfg.GuildID,
		registeredCommands: make([]*discordgo.ApplicationCommand, 0),
		characterComponent: component,
	}, nil
}

func (b *bot) Start() error {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := b.session.Open()

	err = b.characterComponent.Load(b.session)
	if err != nil {
		return err
	}

	err = b.addRonnieDCommand()
	if err != nil {
		log.Print(err)
		return err
	}

	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Content == "ping" {
			msg, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
			if err != nil {
				log.Println(err)
			}
			log.Println(msg)
		}
	})

	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "ronnied":
				b.processRonnieD(s, i)
			}
		}
	})

	return nil
}

func (b *bot) Close() error {

	for _, v := range b.registeredCommands {
		log.Printf("Removing command '%v'...", v.Name)

		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	return b.session.Close()
}
