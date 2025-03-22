package protocol

import (
	"encoding/binary"
	"strconv"
)


type AgencyID struct {
	ID string
}

func NewAgencyID(id string) *AgencyID {
	return &AgencyID{
		ID: id,
	}
}

func (a *AgencyID) Encode() ([]byte, int, error) {
	const ENCODE_LEN = 2
	idInBytes := make([]byte, ENCODE_LEN)

	const BASE = 10
	const BIT_SIZE = 16
	buf, err := strconv.ParseUint(a.ID, BASE, BIT_SIZE)

	if err != nil {
		return nil, 0, err
	}

	binary.BigEndian.PutUint16(idInBytes, uint16(buf))
	
	return idInBytes, ENCODE_LEN, nil
}