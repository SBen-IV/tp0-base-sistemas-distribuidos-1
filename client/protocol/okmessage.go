package protocol

type OKMessage struct {
	Buf []byte
	BytesAmount int
}

type ServerMessage int

// Messages sent by the server
const (
	Ok ServerMessage = iota
	ErrorCode1
	Unkown
)

const (
	okMessage = "OK"
	errorCode1Message = "E1"
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
	case errorCode1Message:
		response = ErrorCode1
	}

	return response
}
