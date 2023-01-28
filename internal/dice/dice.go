package dice

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type RollResult struct {
	Used    bool
	Total   int
	Highest int
	Lowest  int
	Rolls   []int
}

func Roll(count, size int) (*RollResult, error) {
	if count < 1 {
		return nil, errors.New("invalid dice count")
	}

	if size < 1 {
		return nil, errors.New("invalid dice size")
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

	log.Println("Rolling", count, "d", size, ":", out, "total:", total, "min:", min, "max:", max)
	return &RollResult{
		Total:   total,
		Highest: max,
		Lowest:  min,
		Rolls:   out,
	}, nil
}

func RollString(diceString string) (*RollResult, error) {
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

	return Roll(count, size)
}

func (r *RollResult) Display() string {
	return fmt.Sprintf("*%d* : %v  ", r.Total-r.Lowest, r.Rolls)
}
