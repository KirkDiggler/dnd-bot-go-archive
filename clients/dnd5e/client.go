package dnd5e

import (
	"net/http"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
	apiEntities "github.com/fadedpez/dnd5e-api/entities"

	"github.com/KirkDiggler/dnd-bot-go/dnderr"
	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
)

// TODO: add context to functions
type client struct {
	client dnd5e.Interface
}

type Config struct {
	HttpClient *http.Client
}

func New(cfg *Config) (Client, error) {
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

func (c *client) GetRace(key string) (*entities.Race, error) {
	response, err := c.client.GetRace(key)
	if err != nil {
		return nil, err
	}

	race := apiRaceToRace(response)

	return race, nil
}

func (c *client) GetClass(key string) (*entities.Class, error) {
	response, err := c.client.GetClass(key)
	if err != nil {
		return nil, err
	}

	return apiClassToClass(response), nil
}

func (c *client) GetProficiency(key string) (*entities.Proficiency, error) {
	if key == "" {
		return nil, dnderr.NewMissingParameterError("GetProficiency.key")
	}

	response, err := c.doGetProficiency(key)
	if err != nil {
		return nil, err
	}

	return apiProficiencyToProficiency(response), nil
}

func (c *client) doGetProficiency(key string) (*apiEntities.Proficiency, error) {
	response, err := c.client.GetProficiency(key)
	if err != nil {
		return nil, err
	}

	return response, nil
}
