package entities

type ProficiencyChoices struct {
	Selected bool
	Count    int
	From     []*Proficiency
}

type Proficiency struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
