package discordbot

type Interface interface {
	// Start starts the bot
	Start() error
	// Close closes the bot
	Close() error
}
