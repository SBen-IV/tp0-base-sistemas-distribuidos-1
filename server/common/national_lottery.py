import logging

from model.winner import Winner
from common.utils import Bet, has_won, load_bets, store_bets

class SingletonNationalLottery(type):
    _instances = {}

    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            instance = super().__call__(*args, **kwargs)
            cls._instances[cls] = instance
        return cls._instances[cls]


class NationalLottery(metaclass=SingletonNationalLottery):
    """
    Represents the National Lottery. Handles bets storage and draw.
    """
    def __init__(self, clients_amount: int):
        self._agencies = {}
        self._clients_amount = clients_amount

    def add_agency(self, agency_id: str):
        value = self._agencies.get(agency_id, None)

        if value is None:
            self._agencies[agency_id] = False

    def mark_no_more_bets(self, agency_id: str):
        self._agencies[agency_id] = True

    def store_bets(self, bets: list[Bet]):
        store_bets(bets)
            
    def can_draw(self) -> bool:
        all_clients_reported = self._clients_amount == len(self._agencies)
        all_clients_done = all(value for value in self._agencies.values())

        logging.debug(f"can_draw({all_clients_reported}, {all_clients_done})")

        return all_clients_reported and all_clients_done
    
    def get_winners(self, agency_id: str) -> list[Winner]:
        winners = []

        for bet in list(load_bets()):
            if bet.agency == int(agency_id) and has_won(bet):
                winners.append(Winner(bet.document))

        return winners