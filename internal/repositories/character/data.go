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
	Str *AbilityScoreData `json:"str"`
	Dex *AbilityScoreData `json:"dex"`
	Con *AbilityScoreData `json:"con"`
	Int *AbilityScoreData `json:"int"`
	Wis *AbilityScoreData `json:"wis"`
	Cha *AbilityScoreData `json:"cha"`
}

type AbilityScoreData struct {
	Score int `json:"score"`
	Bonus int `json:"bonus"`
}
type Proficiency struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}
