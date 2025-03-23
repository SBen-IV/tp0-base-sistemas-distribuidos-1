from common.utils import Bet
from protocol.betfrombytes_error import BetFromBytesError
from protocol.bet import BetProtocol


class BetsProtocol():
    @staticmethod
    def from_bytes(msg, agency_id: str, bets_amount: int) -> list[Bet]:
        bets = []
        bet_as_str = msg.decode("utf-8")
        bets_split = bet_as_str.split(",")

        # Remove last empty string
        if bets_split[-1] == "":
            bets_split = bets_split[:-1]

        for bet in bets_split:
            bets.append(BetProtocol.from_bytes(bet, agency_id)) 

        if len(bets) < bets_amount:
            raise BetFromBytesError(f"Less bets ({len(bets)}) than expected ({bets_amount})")

        return bets
    