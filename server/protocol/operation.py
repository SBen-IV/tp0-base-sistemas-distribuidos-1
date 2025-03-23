from model.operation import Operation

BET_MESSAGE = "BET"
FIN_MESSAGE = "FIN"

OPERATION_MESSAGE_LEN = 3

class OperationProtocol():
    @staticmethod
    def from_bytes(buf):
        msg = buf.decode('utf-8')

        if msg == BET_MESSAGE:
            return Operation.Bet
        elif msg == FIN_MESSAGE:
            return Operation.Fin
        
        return Operation.Unkown