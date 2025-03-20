package common

import (
	"encoding/binary"
	"strconv"
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