package character

import "github.com/KirkDiggler/dnd-bot-go/internal/entities"

type Data struct {
	ID            string         `json:"id"`
	OwnerID       string         `json:"owner_id"`
	Name          string         `json:"name"`
	ClassKey      string         `json:"class_key"`
	Class         *ClassData     `json:"class"`
	RaceKey       string         `json:"race_key"`
	Attributes    *AttributeData `json:"attributes"`
	Rolls         []*RollData    `json:"rolls"`
	Proficiencies []*Proficiency `json:"proficiencies"`
}
type ClassData struct {
	Key                string             `json:"key"`
	ProficiencyChoices []*entities.Choice `json:"proficiency_choices"`
}

type RollData struct {
	Used    bool  `json:"used"`
	Total   int   `json:"total"`
	Highest int   `json:"highest"`
	Lowest  int   `json:"lowest"`
	Rolls   []int `json:"rolls"`
}

type AttributeData struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
}

type Proficiency struct {
	Key string `json:"key"`
}

type ProficiencyChoice struct {
	Selected bool `json:"selected"`
	Count    int  `json:"count"`
	From     []*Proficiency
}
