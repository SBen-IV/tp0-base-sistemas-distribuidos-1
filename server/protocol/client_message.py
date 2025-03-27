from enum import Enum


OK_CLIENT_MESSAGE_LEN = 2
OK_CLIENT_MESSAGE = "OK"

class ClientResponseMessage(Enum):
    OkMessage = 1
    Unkown = 2


class ClientMessage():
    @staticmethod
    def from_bytes(buf, bytes_amount):
        message = buf[0:bytes_amount].decode("utf-8")

        if message == OK_CLIENT_MESSAGE:
            return ClientResponseMessage.OkMessage
        
        return ClientResponseMessage.Unkown

