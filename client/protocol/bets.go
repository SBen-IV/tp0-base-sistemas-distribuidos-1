package protocol

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"

type BetsProtocol struct {
	bets []Bet
}

func NewBets(bets []model.Bet) BetsProtocol {
	betsProtocol := []Bet{}

	for _, modelBet := range bets {
		bet := NewBet(&modelBet)

		betsProtocol = append(betsProtocol, *bet)
	}

	return BetsProtocol{
		bets: betsProtocol,
	}
}

func (b *BetsProtocol) Encode() ([]byte, int) {
	buf := []byte{}
	totalBytesAmount := 0

	for _, bet := range b.bets {
		betAsBytes, bytesAmount := bet.Encode()

		buf = append(buf, betAsBytes...)
		totalBytesAmount += bytesAmount
	}

	return buf, totalBytesAmount
}