class SingletonNationalLottery(type):
    _instance = None

    def __call__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super().__call__(*args, **kwargs)
        return cls._instance


class NationalLottery(metaclass=SingletonNationalLottery):
    def __init__(self):
        self._agencies = {}

    def add_agency(self, agency_id: str):
        if agency_id not in self._agencies:
            self._agencies[agency_id] = False

    def mark_no_more_bets(self, agency_id: str):
        self._agencies[agency_id] = True
            
    def can_draw(self) -> bool:
        return all(value for value in self._agencies.values())
    