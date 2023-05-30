package entities

type Monster struct {
	ID          string `json:"id"`
	CharacterID string `json:"character_id"`
	CurrentHP   int    `json:"current_hp"`
	Key         string `json:"key"`
}
