from common.utils import Bet
from protocol.betfrombytes_error import BetFromBytesError

FIRST_NAME_POS = 0
LAST_NAME_POS = 1
DOCUMENT_POS = 2
BIRTHDAY_POS = 3
NUMBER_POS = 4
COMPONENTS_LEN = 5

class BetProtocol():
    @staticmethod
    def from_bytes(msg, agency_id) -> Bet:
        bet_as_str = msg.decode("utf-8")
        bet_components = bet_as_str.split(";")

        if len(bet_components) < COMPONENTS_LEN:
            raise BetFromBytesError(f"Less components ({len(bet_components)}) than expected ({COMPONENTS_LEN})")

        return Bet(agency=agency_id,
            first_name=bet_components[FIRST_NAME_POS],
            last_name=bet_components[LAST_NAME_POS],
            document=bet_components[DOCUMENT_POS],
            birthdate=bet_components[BIRTHDAY_POS],
            number=bet_components[NUMBER_POS]
        )