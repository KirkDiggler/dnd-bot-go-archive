package entities

type Proficiency struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (p *Proficiency) String() string {
	return p.Name
}
