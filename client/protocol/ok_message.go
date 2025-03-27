package protocol

type OkMessage struct {
	data string
}

func NewOkMessageBuf() OkMessage {
	return OkMessage{data: "OK"}
}

func (o *OkMessage) Encode() ([]byte, int) {
	buf := []byte(o.data)

	return buf, len(buf)
}