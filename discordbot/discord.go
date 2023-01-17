package discordbot

import (
	"github.com/KirkDiggler/dnd-bot-go/errors"
	"github.com/bwmarrin/discordgo"
	"log"
)

type impl struct {
	bot *discordgo.Session
}

type Config struct {
	Token string
}

func New(cfg *Config) (*impl, error) {
	if cfg == nil {
		return nil, errors.NewMissingParameterError("cfg")
	}
	if cfg.Token == "" {
		return nil, errors.NewMissingParameterError("cfg.Token")
	}

	bot, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Content == "ping" {
			msg, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
			if err != nil {
				log.Println(err)
			}
			log.Println(msg)
		}
	})

	return &impl{
		bot: bot,
	}, nil
}

func (i *impl) Start() error {
	return i.bot.Open()
}

func (i *impl) Close() error {
	return i.bot.Close()
}

func (i *impl) ronnieD(handler interface{}) {
	i.bot.AddHandler(handler)
}
