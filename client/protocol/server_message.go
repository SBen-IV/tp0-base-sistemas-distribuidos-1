package protocol

type ServerMessage int

// Messages sent by the server
const (
	Ok ServerMessage = iota
	NotOk
	Unkown
)

const (
	okMessage = "OK"
	notOkMessage = "NO"
)

func NewMessageBuf() ([]byte, int) {
	const MSG_BUF_LEN = 2

	return make([]byte, MSG_BUF_LEN), MSG_BUF_LEN
}

func NewServerMessageBuild(buf []byte, bytesAmount int) ServerMessage {
	var response ServerMessage = Unkown

	switch string(buf[0:bytesAmount]) {
	case okMessage:
		response = Ok
	case notOkMessage:
		response = NotOk
	}

	return response
}
