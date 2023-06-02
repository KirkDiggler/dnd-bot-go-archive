package damage

import (
	"errors"
	"github.com/KirkDiggler/dnd-bot-go/internal/dice"
	"strconv"
	"strings"
)

type Type string

const (
	TypeAcid        Type = "acid"
	TypeCold        Type = "cold"
	TypeFire        Type = "fire"
	TypeForce       Type = "force"
	TypeLightning   Type = "lightning"
	TypeNecrotic    Type = "necrotic"
	TypePoison      Type = "poison"
	TypePsychic     Type = "psychic"
	TypeRadiant     Type = "radiant"
	TypeThunder     Type = "thunder"
	TypeBludgeoning Type = "bludgeoning"
	TypePiercing    Type = "piercing"
	TypeSlashing    Type = "slashing"
	TypeNone        Type = "none"
)

type Damage struct {
	DiceCount  int
	DiceSize   int
	Bonus      int
	DamageType Type
}

func (d *Damage) Deal() int {
	base, _ := dice.Roll(d.DiceCount, d.DiceSize, d.Bonus)
	return base.Total
}

func DamageDiceToDamage(input string) (*Damage, error) {
	a := strings.Split(input, "+")
	var diceString = input
	var bonus, diceSize, diceCount int
	var err error
	if len(a) == 2 {
		bonus, _ = strconv.Atoi(a[1])
		diceString = a[0]
	}

	diceParts := strings.Split(diceString, "d")
	if len(diceParts) != 2 {
		return nil, errors.New("invalid dice string")
	}

	strCount := diceParts[0]
	strSize := diceParts[1]

	diceCount, err = strconv.Atoi(strCount)
	if err != nil {
		return nil, errors.New("invalid dice string")
	}
	diceSize, err = strconv.Atoi(strSize)
	if err != nil {
		return nil, errors.New("invalid dice string")
	}

	return &Damage{
		DiceCount: diceCount,
		DiceSize:  diceSize,
		Bonus:     bonus,
	}, nil
}
