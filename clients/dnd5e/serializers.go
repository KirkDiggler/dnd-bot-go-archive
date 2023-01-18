package dnd5e

import (
	"github.com/KirkDiggler/dnd-bot-go/entities"
	apiEntities "github.com/fadedpez/dnd5e-api/entities"
)

func apiReferenceItemToClass(apiClass *apiEntities.ReferenceItem) *entities.Class {
	return &entities.Class{
		Key:  apiClass.Key,
		Name: apiClass.Name,
	}
}

func apiReferenceItemsToClasses(input []*apiEntities.ReferenceItem) []*entities.Class {
	output := make([]*entities.Class, len(input))
	for i, apiClass := range input {
		output[i] = apiReferenceItemToClass(apiClass)
	}
	return output
}

func apiReferenceItemToRace(input *apiEntities.ReferenceItem) *entities.Race {
	return &entities.Race{
		Key:  input.Key,
		Name: input.Name,
	}
}

func apiReferenceItemsToRaces(input []*apiEntities.ReferenceItem) []*entities.Race {
	output := make([]*entities.Race, len(input))
	for i, apiRace := range input {
		output[i] = apiReferenceItemToRace(apiRace)
	}

	return output
}
