from model.winner import Winner


class WinnersMessageProtocol():
    @staticmethod
    def to_bytes(winners: list[Winner]):
        buf = b''

        for winner in winners:
            buf += winner._document.encode('utf-8') + b';'

        return buf, len(buf)