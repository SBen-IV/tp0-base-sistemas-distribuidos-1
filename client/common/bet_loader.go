package common

import (
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

type BetLoader interface {
	GetBet() *model.Bet
}

// Reads bet from env variables
type betLoader struct {}

func CreateBetLoader() *betLoader {
	return &betLoader{}
}

func (b *betLoader) GetBet() *model.Bet {
	firstName := os.Getenv("NOMBRE")
	lastName := os.Getenv("APELLIDO")
	document := os.Getenv("DOCUMENTO")
	birthday := os.Getenv("NACIMIENTO")
	numberStr := os.Getenv("NUMERO")

	number, err := strconv.ParseInt(numberStr, 10, 32)

	if err != nil {
		// Should never return error
		return nil
	}

	return model.NewBet(firstName, lastName, document, birthday, int32(number))
}
