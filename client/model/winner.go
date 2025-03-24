package model

type Winner struct {
	Document string
}

func NewWinner(document string) Winner {
	return Winner{Document: document}
}