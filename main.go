package main

import (
	"flag"
	"github.com/KirkDiggler/dnd-bot-go/discordbot"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	token = flag.String("token", "", "Bot token")
)

func main() {
	flag.Parse()
	if token == nil {
		panic("Token is required")
	}

	if *token == "" {
		flag.Usage()
		return
	}

	bot, err := discordbot.New(&discordbot.Config{
		Token: *token,
	})
	if err != nil {
		panic(err)
	}

	err = bot.Start()
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
end:
	for {
		select {
		case <-stchan:
			break end
		default:
		}
		time.Sleep(time.Second)
	}

}
