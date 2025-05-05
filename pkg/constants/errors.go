package constants

import "errors"

var (
	ErrInvalidUserID  = errors.New("invalid user ID")
	ErrInvalidEventID = errors.New("invalid event ID")

	ErrInvalidEventStatus       = errors.New("invalid event status")
	ErrInvalidEventResultStatus = errors.New("invalid event result status")
	ErrInvalidOdds              = errors.New("invalid bet odds")
	ErrInvalidAmount            = errors.New("invalid bet amount")

	ErrInsuffiecientBalance = errors.New("insuffiecient balance")
	ErrBetAlreadyResulted   = errors.New("bet already resulted")
	ErrEventAlreadyResulted = errors.New("event already resulted")
)
