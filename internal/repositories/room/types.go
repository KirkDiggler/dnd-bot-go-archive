package room

type ListByPlayerInput struct {
	PlayerID string
	Limit    int64
	Offset   int64
	Reverse  bool
}
