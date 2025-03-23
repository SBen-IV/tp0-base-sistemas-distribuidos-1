package common

import (
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

type BetLoader interface {
	Init() error
	GetBet() *model.Bet
	Destroy()
}

// Reads bet from env variables
type betLoader struct {
	filename string
	file *os.File
}

func CreateBetLoader(filename string) *betLoader {
	return &betLoader{
		filename: filename,
		file: nil,
	}
}

func (b *betLoader) Init() error {
	file, err := os.Open(b.filename)

	if err != nil {
		return err
	}

	b.file = file

	return nil
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

func (b *betLoader) Destroy() {
	if b.file != nil {
		b.file.Close()
	}
}