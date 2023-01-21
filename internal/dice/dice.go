package dice

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

type RollResult struct {
	Total   int
	Highest int
	Lowest  int
	Details []int
}

func Roll(diceString string) (*RollResult, error) {
	diceParts := strings.Split(diceString, "d")
	if len(diceParts) != 2 {
		return nil, errors.New("invalid dice string")
	}

	strCount := diceParts[0]
	strSize := diceParts[1]

	count, err := strconv.Atoi(strCount)
	if err != nil {
		return nil, errors.New("invalid dice string")
	}
	size, err := strconv.Atoi(strSize)
	if err != nil {
		return nil, errors.New("invalid dice string")
	}
	max, min, total := 0, 0, 0

	out := make([]int, count)
	for i := 0; i < count; i++ {
		roll := rand.Intn(size) + 1
		total += roll
		if i == 0 {
			min = roll
			max = roll
		}

		if min > roll {
			min = roll
		}

		if max < roll {
			max = roll
		}

		out[i] = roll
	}

	return &RollResult{
		Total:   total,
		Highest: max,
		Lowest:  min,
		Details: out,
	}, nil
}
