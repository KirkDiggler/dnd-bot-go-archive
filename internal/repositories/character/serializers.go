package character

import (
	"encoding/json"
)

func jsonToData(input string) *Data {
	char := &Data{}

	err := json.Unmarshal([]byte(input), char)
	if err != nil {
		return nil
	}

	return char
}

func dataToJSON(input *Data) string {
	b, err := json.Marshal(input)
	if err != nil {
		return ""
	}

	return string(b)
}
