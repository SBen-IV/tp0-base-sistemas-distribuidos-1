import logging
from threading import Lock


class AgenciesSafe():
    def __init__(self, agencies_amount: int):
        self._agencies = {}
        self._lock = Lock()
        self._agencies_amount = agencies_amount

    def add_agency(self, agency_id: str):
        with self._lock:
            value = self._agencies.get(agency_id, None)

            if value is None:
                self._agencies[agency_id] = False
    
    def mark_no_more_bets(self, agency_id: str):
        with self._lock:
            self._agencies[agency_id] = True

    def can_draw(self) -> bool:
        are_equal_amount = False
        all_agencies_done = False

        with self._lock:
            are_equal_amount = self._agencies_amount == len(self._agencies)
            all_agencies_done = all(v for v in self._agencies.values())
            
        logging.debug(f"can_draw: {are_equal_amount} and {all_agencies_done}")

        return are_equal_amount and all_agencies_done