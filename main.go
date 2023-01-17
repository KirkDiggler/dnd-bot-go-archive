package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KirkDiggler/dnd-bot-go/discordbot"
)

var (
	token   string
	guildID string
)

func init() {
	flag.StringVar(&token, "t", "",
		"Bot token")
	flag.StringVar(&guildID, "g", "",
		"Guild ID")
}

func main() {
	flag.Parse()

	if token == "" {
		flag.Usage()
		return
	}

	if guildID == "" {
		flag.Usage()
		return
	}

	bot, err := discordbot.New(&discordbot.Config{
		Token:   token,
		GuildID: guildID,
	})
	if err != nil {
		panic(err)
	}

	err = bot.Start()
	if err != nil {
		panic(err)
	}

	defer func(bot discordbot.Interface) {
		err := bot.Close()
		if err != nil {
			panic(err)
		}
	}(bot)

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)

	for {
		select {
		case <-stchan:
			return
		default:
		}
		time.Sleep(time.Second)
	}

}
