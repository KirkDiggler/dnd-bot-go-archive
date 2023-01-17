package discordbot

import (
	"log"

	"github.com/KirkDiggler/dnd-bot-go/errors"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session *discordgo.Session
	guildID string

	registeredCommands []*discordgo.ApplicationCommand
}

type Config struct {
	Token   string
	GuildID string
}

func New(cfg *Config) (*bot, error) {
	if cfg == nil {
		return nil, errors.NewMissingParameterError("cfg")
	}

	if cfg.Token == "" {
		return nil, errors.NewMissingParameterError("cfg.Token")
	}

	if cfg.GuildID == "" {
		return nil, errors.NewMissingParameterError("cfg.GuildID")
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	return &bot{
		session:            session,
		guildID:            cfg.GuildID,
		registeredCommands: make([]*discordgo.ApplicationCommand, 0),
	}, nil
}

func (b *bot) Start() error {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := b.session.Open()

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

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
