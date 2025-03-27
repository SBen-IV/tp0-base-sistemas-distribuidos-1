package protocol

import "encoding/binary"

// Information about the bet sent to the server before sending the bet
type BetInfo struct {
	BetsAmount int32
	BytesAmount int32
}

func NewBetInfo(bets_amount int32, bytes_amount int32) *BetInfo {
	return &BetInfo{
		BetsAmount: bets_amount,
		BytesAmount: bytes_amount,
	}
}

func (b *BetInfo) Encode() ([]byte, int, error) {
	const ENCODE_LEN = 8

	buf := make([]byte, ENCODE_LEN)

	binary.BigEndian.PutUint32(buf[0:4], uint32(b.BetsAmount))
	binary.BigEndian.PutUint32(buf[4:8], uint32(b.BytesAmount))

	return buf, ENCODE_LEN, nil
}