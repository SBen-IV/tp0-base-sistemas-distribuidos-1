from protocol.response_message import ResponseMessage


class OkMessage(ResponseMessage):
    def __init__(self):
        super().__init__("OK")