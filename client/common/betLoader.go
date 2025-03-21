package common

import (
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

type BetLoader interface {
	Init() error
	GetBet() *model.Bet
	Destroy() error
}

type betLoader struct {}

func CreateBetLoader() *betLoader {
	return &betLoader{}
}

func (b *betLoader) Init() error {
	return nil
}

func (b *betLoader) GetBet() *model.Bet {
	first_name := os.Getenv("NOMBRE")
	last_name := os.Getenv("APELLIDO")
	document := os.Getenv("DOCUMENTO")
	birthday := os.Getenv("NACIMIENTO")
	number_str := os.Getenv("NUMERO")

	number, err := strconv.ParseInt(number_str, 10, 32)

	if err != nil {
		// Should never return error
		return nil
	}

	return model.NewBet(first_name, last_name, document, birthday, int32(number))
	
}

func (b *betLoader) Destroy() error {
	return nil
}