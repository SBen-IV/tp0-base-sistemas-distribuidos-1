package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

type BetLoader interface {
	Init() error
	GetBet() *model.Bet
	GetBets(maxAmount int) ([]model.Bet, error)
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

func (b *betLoader) GetBets(maxAmount int) ([]model.Bet, error) {
	bets := []model.Bet{}

	scanner := bufio.NewScanner(b.file)

	for i := 0; i < maxAmount; i++ {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			
			break
		}

		// Example
		// Santiago Lionel,Lorca,30904465,1999-03-17,2201
		// line := scanner.Text()
		// parsedLine, err := parseLine(line)

		// if err != nil {
		// 	return nil, err
		// }

		// bet := getBet(parsedLine)
		// bet := strings.Split(, ",")
		// bets = append(bets, bet)
	}



	return bets, nil
}

func (b *betLoader) Destroy() {
	if b.file != nil {
		b.file.Close()
	}
}

func parseLine(line string) ([]string, error) {
	parsedLine := strings.Split(line, ",")

	if len(parsedLine) < 5 {
		return nil, fmt.Errorf("Malformed line")
	}

	return parsedLine, nil
}
