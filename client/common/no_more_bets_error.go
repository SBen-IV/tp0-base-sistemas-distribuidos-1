package common

import "fmt"

type NoMoreBets struct {}

func (e *NoMoreBets) Error() string {
	return fmt.Sprint("No more bets")
}