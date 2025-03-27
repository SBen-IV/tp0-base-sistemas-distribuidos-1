from protocol.response_message import ResponseMessage


class ErrMessage(ResponseMessage):
    def __init__(self):
        super().__init__("NO")
    