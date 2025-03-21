package model

type Bet struct {
	FirstName string
	LastName string
	Document string
	Birthday string
	Number int32
}


func NewBet(firstName string, lastName string, document string, birthday string, number int32) *Bet {
	return &Bet{
		FirstName: firstName,
		LastName: lastName,
		Document: document,
		Birthday: birthday,
		Number: number,
	}
}
