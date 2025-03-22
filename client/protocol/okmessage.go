package protocol

type OKMessage struct {
	Buf []byte
	BytesAmount int
}

type ServerMessage int

const (
	Ok ServerMessage = iota
	Unkown
)

const (
	okMessage = "OK"
)

func NewMessageBuf() ([]byte, int) {
	const MSG_BUF_LEN = 2

	return make([]byte, MSG_BUF_LEN), MSG_BUF_LEN
}

func NewOKMessage() OKMessage {
	const OK_MSG_LEN = 2

	return OKMessage{
		Buf: make([]byte, OK_MSG_LEN),
		BytesAmount: OK_MSG_LEN,
	}
}

func NewServerMessageBuild(buf []byte, bytesAmount int) ServerMessage {
	message := string(buf[0:bytesAmount])

	if message == okMessage {
		return Ok
	}

	return Unkown
}
