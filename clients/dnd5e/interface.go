package dnd5e

import "github.com/KirkDiggler/dnd-bot-go/entities"

type Interface interface {
	ListClasses() ([]*entities.Class, error)
	ListRaces() ([]*entities.Race, error)
}
