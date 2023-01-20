package party

import (
	"encoding/json"

	"github.com/KirkDiggler/dnd-bot-go/entities"
)

func jsonToParty(input string) *entities.Party {
	party := entities.Party{}

	err := json.Unmarshal([]byte(input), &party)
	if err != nil {
		return nil
	}

	return &party
}

func partyToJson(input *entities.Party) string {
	buf, _ := json.Marshal(input)

	return string(buf)
}
