package protocol

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"

type Operation struct {
	data string
}

const (
	BET_OPERATION = "BET"
	DRAW_OPERATION = "DRW"
	NO_MORE_BETS_OPERATION = "NMB"
	FIN_OPERATION = "FIN"
)

func NewOperation(op model.Operation) *Operation {
	data := FIN_OPERATION

	switch op {
	case model.BetOp:
		data = BET_OPERATION
	case model.DrawOp:
		data = DRAW_OPERATION
	case model.NoMoreBets:
		data = NO_MORE_BETS_OPERATION
	}

	return &Operation{data: data}
}

func (o *Operation) Encode() ([]byte, int) {
	buf := []byte(o.data)

	return buf, len(buf)
}