package common

type NoMoreBets struct {}

func (e *NoMoreBets) Error() string {
	return "no more bets"
}