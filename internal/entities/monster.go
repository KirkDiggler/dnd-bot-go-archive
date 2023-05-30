package entities

type Monster struct {
	ID         string `json:"id"`
	CharcterID string `json:"character_id"`
	CurrentHP  int    `json:"current_hp"`
	Key        string `json:"key"`
}
