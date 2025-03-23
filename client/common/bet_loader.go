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

		line := scanner.Text()
		parsedLine, err := b.parseLine(line)

		if err != nil {
			return nil, err
		}

		bet, err := b.getBet(parsedLine)

		if err != nil {
			return nil, err
		}
		
		bets = append(bets, *bet)
	}

	return bets, nil
}

func (b *betLoader) Destroy() {
	if b.file != nil {
		b.file.Close()
	}
}

func (b *betLoader) parseLine(line string) ([]string, error) {
	const COMPONENTS_LEN = 5
	const COMPONENT_SEPARATOR = ","
	// Example line
	// Santiago Lionel,Lorca,30904465,1999-03-17,2201
	parsedLine := strings.Split(line, COMPONENT_SEPARATOR)

	if len(parsedLine) < COMPONENTS_LEN {
		return nil, fmt.Errorf("malformed line")
	}

	return parsedLine, nil
}

func (b *betLoader) getBet(parsedLine []string) (*model.Bet, error) {
	const FIRST_NAME_POS = 0
	const LAST_NAME_POS = 1
	const DOCUMENT_POS = 2
	const BIRTHDAY_POS = 3
	const NUMBER_POS = 4

	firstName := parsedLine[FIRST_NAME_POS]
	lastName := parsedLine[LAST_NAME_POS]
	document := parsedLine[DOCUMENT_POS]
	birthday := parsedLine[BIRTHDAY_POS]
	numberStr := parsedLine[NUMBER_POS]

	number, err := strconv.ParseInt(numberStr, 10, 32)

	if err != nil {
		// Should never return error
		return nil, err
	}

	return model.NewBet(firstName, lastName, document, birthday, int32(number)), nil
}