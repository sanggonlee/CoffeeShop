package utils

import "errors"

var (
	ErrMoneyFormat      = errors.New("Unrecognized format for money")
	ErrMoneyDollarParse = errors.New("Failed parsing dollar part")
	ErrMoneyCentParse   = errors.New("Failed parsing cent part")
	ErrMoneyAmount      = errors.New("Money given is out of accepted range")

	ErrStartEndDate = errors.New("The end date comes before the start date")

	ErrDeleteMissingID = errors.New("ID is required for delete operation")
)
