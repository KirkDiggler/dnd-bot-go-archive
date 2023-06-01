package combat_log

type ListByEncounterInput struct {
	EncounterID string
	Limit       int
	Offset      int
	Reverse     bool
}
