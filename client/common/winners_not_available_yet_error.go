package common

type WinnersNotAvailableYet struct {}

func (e *WinnersNotAvailableYet) Error() string {
	return "winners not available yet"
}