class OkMessage():
    def __init__(self):
        self._data = "OK"

    def encode(self):
        return self._data.encode("utf-8")