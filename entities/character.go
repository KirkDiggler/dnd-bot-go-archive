package entities

type Character struct {
	ID    string `json:"id"`
	Name  string
	Race  *Race
	Class *Class
}
