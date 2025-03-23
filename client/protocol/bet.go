package protocol

import (
	"fmt"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

// Bet used for protocol layer
type Bet struct {
	FirstName string
	LastName string
	Document string
	Birthday string
	Number string
}

func NewBet(bet *model.Bet) *Bet {
	number := strconv.Itoa(int(bet.Number))

	return &Bet{
		FirstName: bet.FirstName,
		LastName: bet.LastName,
		Document: bet.Document,
		Birthday: bet.Birthday,
		Number: number,
	}
}

func (b *Bet) Encode() ([]byte, int) {
	message := fmt.Sprintf("%s;%s;%s;%s;%s,", b.FirstName, b.LastName, b.Document, b.Birthday, b.Number)

	buf := []byte(message)

	return buf, len(buf)
}