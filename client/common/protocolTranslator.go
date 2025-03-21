package common

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

type ProtocolTranslator struct {}


func NewProtocolTranslator() *ProtocolTranslator {
	return &ProtocolTranslator{}
}

func (p *ProtocolTranslator) IDtoBytes(id string) ([]byte, error) {
	const BASE = 10
	const BIT_SIZE = 16
	message_id, err := strconv.ParseUint(id, BASE, BIT_SIZE)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 2)

	binary.BigEndian.PutUint16(buf, uint16(message_id))

	return buf, nil
}

func (p *ProtocolTranslator) BetToBytes(bet *model.Bet) []byte {
	message := fmt.Sprintf("%s;%s;%s;%s;%d;", bet.FirstName, bet.LastName, bet.Document, bet.Birthday, bet.Number)

	buf := []byte(message)

	return buf
}

func (p *ProtocolTranslator) OKtoString(buf []byte) (string, error) {
	return string(buf), nil
}