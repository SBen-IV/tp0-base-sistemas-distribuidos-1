BUF_LEN = 4

class WinnersInfoMessageProtocol():
    @staticmethod
    def to_bytes(bytes_amount: int):
        buf = int.to_bytes(bytes_amount, BUF_LEN, "big")

        return buf, BUF_LEN