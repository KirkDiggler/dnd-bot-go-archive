package dnd5e

import (
	"net/http"

	"github.com/KirkDiggler/dnd-bot-go/entities"

	"github.com/KirkDiggler/dnd-bot-go/errors"
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
		return nil, errors.NewMissingParameterError("cfg")
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

func (c *client) ListClasses() ([]*entities.Class, error) {
	response, err := c.client.ListClasses()
	if err != nil {
		return nil, err
	}

	return apiReferenceItemsToClasses(response), nil
}

func (c *client) ListRaces() ([]*entities.Race, error) {
	response, err := c.client.ListRaces()
	if err != nil {
		return nil, err
	}

	return apiReferenceItemsToRaces(response), nil
}
