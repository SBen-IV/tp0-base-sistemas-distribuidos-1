package protocol

import (
	"encoding/binary"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"
)

const BUF_LEN = 4

func NewWinnersBytesAmountBuf() ([]byte, int) {
	buf := make([]byte, BUF_LEN)

	return buf, BUF_LEN
}

func NewWinnersBytesAmountBuild(buf []byte) int {
	bytesAmount := binary.BigEndian.Uint32(buf[0:BUF_LEN])

	return int(bytesAmount)
}

func NewWinnersBuf(bytesAmount int) ([]byte, int) {
	return make([]byte, bytesAmount), bytesAmount
}

func NewWinnersBuild(buf []byte, bytesAmount int) ([]model.Winner, error) {
	message := string(buf[0:bytesAmount])
	documents := strings.Split(message, ";")

	// Remove empty element
	if documents[len(documents)-1] == "" {
		documents = documents[:len(documents)-1]
	}

	winners := []model.Winner{}

	for _, winner := range documents {
		winners = append(winners, model.NewWinner(winner))
	}

	return winners, nil
}