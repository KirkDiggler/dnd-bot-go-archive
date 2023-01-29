package character

type Data struct {
	ID            string         `json:"id"`
	OwnerID       string         `json:"owner_id"`
	Name          string         `json:"name"`
	ClassKey      string         `json:"class_key"`
	RaceKey       string         `json:"race_key"`
	Attributes    *AttributeData `json:"attributes"`
	Rolls         []*RollData    `json:"rolls"`
	Proficiencies []*Proficiency `json:"proficiencies"`
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
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}
