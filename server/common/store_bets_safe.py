from threading import Lock

from common.utils import Bet, load_bets, store_bets


class BetsStorageSafe():
    def __init__(self):
        self._lock = Lock()

    def store_bets(self, bets: list[Bet]):
        with self._lock:
            store_bets(bets)

    def load_bets(self):
        with self._lock:
            return load_bets()