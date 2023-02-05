package character

type Data struct {
	ID               string         `json:"id"`
	OwnerID          string         `json:"owner_id"`
	Name             string         `json:"name"`
	ClassKey         string         `json:"class_key"`
	RaceKey          string         `json:"race_key"`
	AC               int            `json:"ac"`
	Speed            int            `json:"speed"`
	HitDie           int            `json:"hit_die"`
	Level            int            `json:"level"`
	Experience       int            `json:"experience"`
	MaxHitPoints     int            `json:"max_hit_points"`
	CurrentHitPoints int            `json:"current_hit_points"`
	Attributes       *AttributeData `json:"attributes"`
	NextLevel        int            `json:"next_level"`
	Rolls            []*RollData    `json:"rolls"`
	Proficiencies    []*Proficiency `json:"proficiencies"`
	Inventory        []*Equipment   `json:"inventory"`
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
	Type string `json:"type,omitempty"`
}

type Equipment struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
