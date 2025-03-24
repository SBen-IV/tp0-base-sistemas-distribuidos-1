from model.winner import Winner
from common.utils import has_won, load_bets


class SingletonNationalLottery(type):
    _instances = {}

    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            instance = super().__call__(*args, **kwargs)
            cls._instances[cls] = instance
        return cls._instances[cls]


class NationalLottery(metaclass=SingletonNationalLottery):
    def __init__(self):
        self._agencies = {}

    def add_agency(self, agency_id: str):
        value = self._agencies.get(agency_id, None)

        if value is None:
            self._agencies[agency_id] = False

    def mark_no_more_bets(self, agency_id: str):
        self._agencies[agency_id] = True
            
    def can_draw(self) -> bool:
        return all(value for value in self._agencies.values())
    
    def get_winners(self, agency_id: str) -> list[Winner]:
        winners = []

        for bet in load_bets():
            if bet.agency == agency_id and has_won(bet):
                winners.append(Winner(bet.document))

        return winners