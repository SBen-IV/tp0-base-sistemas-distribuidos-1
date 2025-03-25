from model.winner import Winner
from common.utils import Bet, has_won, load_bets, store_bets
from server.common.agencies_safe import AgenciesSafe
from server.common.bets_storage_safe import BetsStorageSafe

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
        self._agencies = AgenciesSafe(clients_amount)
        self._bets_storage = BetsStorageSafe()

    def add_agency(self, agency_id: str):
        self._agencies.add_agency(agency_id)

    def mark_no_more_bets(self, agency_id: str):
        self._agencies.mark_no_more_bets(agency_id)

    def store_bets(self, bets: list[Bet]):
        self._bets_storage.store_bets(bets)
            
    def can_draw(self) -> bool:
        return self._agencies.can_draw()
    
    def get_winners(self, agency_id: str) -> list[Winner]:
        winners = []
        bets = list(self._bets_storage.load_bets())

        for bet in bets:
            if bet.agency == int(agency_id) and has_won(bet):
                winners.append(Winner(bet.document))

        return winners