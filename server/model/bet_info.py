class BetInfo():
    """
    Model class that contains information about the incoming Bet
    """
    def __init__(self, bet_amount, bytes_amount):
        self._bets_amount = bet_amount
        self._bytes_amount = bytes_amount