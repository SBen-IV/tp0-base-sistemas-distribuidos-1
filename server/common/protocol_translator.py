from common.bet_info import BetInfo
from common.utils import Bet


class ProtocolTranslator():
    def translate_cli_id(self, cli_id):
        return str(int.from_bytes(cli_id[0:2], 'big'))
    
    def translate_bet_info(self, bet_info) -> BetInfo:
        bet_amount = int.from_bytes(bet_info[0:4], "big")
        bytes_amount = int.from_bytes(bet_info[4:8], "big")

        return BetInfo(bet_amount, bytes_amount)
    
    def translate_bet(self, bet, agency_id) -> Bet:
        bet_as_str = bet.decode('utf-8')
        bet_components = bet_as_str.split(';')

        return Bet(agency=agency_id, first_name=bet_components[0], last_name=bet_components[1], document=bet_components[2], birthdate=bet_components[3], number=bet_components[4])