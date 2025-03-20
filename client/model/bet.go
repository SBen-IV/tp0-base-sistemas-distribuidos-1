package model

type Bet struct {
	firstName string
	lastName string
	document string
	birthday string
	number int32
}


func NewBet(firstName string, lastName string, document string, birthday string, number int32) *Bet {
	return &Bet{
		firstName: firstName,
		lastName: lastName,
		document: document,
		birthday: birthday,
		number: number,
	}
}
