package protocol

import "github.com/7574-sistemas-distribuidos/docker-compose-init/client/model"

type Operation struct {
	data string
}

const (
	BET_OPERATION = "BET"
	DEFAULT_OPERATION = "FIN"
)

func NewOperation(op model.Operation) *Operation {
	data := DEFAULT_OPERATION

	if op == model.BetOp {
		data = BET_OPERATION
	}

	return &Operation{data: data}
}

func (o *Operation) Encode() ([]byte, int) {
	buf := []byte(o.data)

	return buf, len(buf)
}