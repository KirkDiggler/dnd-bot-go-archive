package entities

type Room struct {
	ID          string `json:"id"`
	CharacterID string `json:"character_id"`
	MonsterID   string `json:"monster_id"`
}
