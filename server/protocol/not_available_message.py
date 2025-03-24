from protocol.response_message import ResponseMessage


class NotAvailableMessage(ResponseMessage):
    def __init__(self):
        super().__init__("NA")
