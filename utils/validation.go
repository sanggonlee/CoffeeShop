package utils

import (
	"strconv"
	"strings"
)

const maxDollarAmount = 1000 // completely arbitrary

// ValidateMoney parses money value and performs a sanity check
func ValidateMoney(val string) error {
	tokens := strings.Split(val, ".")
	if len(tokens) != 2 {
		return ErrMoneyFormat
	}

	dollar, err := strconv.ParseInt(tokens[0], 0, 64)
	if err != nil {
		return ErrMoneyDollarParse
	}

	if dollar < 0 || dollar > maxDollarAmount {
		return ErrMoneyAmount
	}

	cent, err := strconv.ParseInt(tokens[1], 0, 64)
	if err != nil {
		return ErrMoneyCentParse
	}

	if cent < 0 || cent > 99 {
		return ErrMoneyCentParse
	}

	return nil
}
