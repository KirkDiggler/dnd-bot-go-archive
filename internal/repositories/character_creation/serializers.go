package character_creation

import (
	"encoding/json"

	"github.com/KirkDiggler/dnd-bot-go/internal/entities"
)

func jsonToCharacterCreation(input string) *entities.CharacterCreation {
	out := &entities.CharacterCreation{}

	err := json.Unmarshal([]byte(input), out)
	if err != nil {
		return nil
	}

	return out
}

func characterCreateToJSON(input *entities.CharacterCreation) string {
	out, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(out)
}
