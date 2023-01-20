package dnd5e

import (
	"net/http"

	entities2 "github.com/KirkDiggler/dnd-bot-go/internal/entities"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
)

type client struct {
	client dnd5e.Interface
}

type Config struct {
	HttpClient *http.Client
}

func New(cfg *Config) (Interface, error) {
	if cfg == nil {
		return nil, dnderr.NewMissingParameterError("cfg")
	}

	dndClient, err := dnd5e.NewDND5eAPI(&dnd5e.DND5eAPIConfig{
		Client: cfg.HttpClient,
	})
	if err != nil {
		return nil, err
	}

	return &client{
		client: dndClient,
	}, nil
}

func (c *client) ListClasses() ([]*entities2.Class, error) {
	response, err := c.client.ListClasses()
	if err != nil {
		return nil, err
	}

	return apiReferenceItemsToClasses(response), nil
}

func (c *client) ListRaces() ([]*entities2.Race, error) {
	response, err := c.client.ListRaces()
	if err != nil {
		return nil, err
	}

	return apiReferenceItemsToRaces(response), nil
}
