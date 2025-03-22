from common.bet_info import BetInfo

BET_INFO_MESSAGE_LEN = 8

class BetInfoProtocol():
    @staticmethod
    def from_bytes(msg):
        bet_amount = int.from_bytes(msg[0:4], "big")
        bytes_amount = int.from_bytes(msg[4:BET_INFO_MESSAGE_LEN], "big")

        return BetInfo(bet_amount, bytes_amount)