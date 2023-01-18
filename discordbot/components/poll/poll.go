package poll

type Poll struct {
	userVotes map[string]int
	votes     map[int]int
}

func New() *Poll {
	return &Poll{
		userVotes: make(map[string]int),
		votes:     make(map[int]int),
	}
}

func (p *Poll) Vote(user string, vote int) {
	if p.userVotes[user] != 0 {
		p.votes[p.userVotes[user]]--
	}
	p.userVotes[user] = vote
	p.votes[vote]++
}

func (p *Poll) GetVotes() map[int]int {
	return p.votes
}
