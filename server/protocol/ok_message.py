class OkMessage():
    def __init__(self):
        self._data = "OK"

    def encode(self):
        data_as_bytes = self._data.encode("utf-8")
        data_len = len(data_as_bytes)

        return data_as_bytes, data_len