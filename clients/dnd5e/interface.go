package dnd5e

import (
	entities2 "github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

type Interface interface {
	ListClasses() ([]*entities2.Class, error)
	ListRaces() ([]*entities2.Race, error)
}
