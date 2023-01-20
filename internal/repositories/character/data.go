package character

type Data struct {
	ID       string `json:"id"`
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	ClassKey string `json:"class_key"`
	RaceKey  string `json:"race_key"`
}
