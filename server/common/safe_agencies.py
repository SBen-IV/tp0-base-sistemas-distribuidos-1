from threading import Lock


class SafeAgencies():
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
        with self._lock:
            return self._agencies_amount == len(self._agencies) and all(v for v in self._agencies.values())